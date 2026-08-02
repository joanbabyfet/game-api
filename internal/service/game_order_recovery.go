package service

import (
	"context"
	"errors"
	"game-api/internal/client/skynet"
	"game-api/internal/model"
	"game-api/pkg"
	"game-api/proto/slotpb"

	"gorm.io/gorm"
)

// RecoverOrder 根据订单创建时的钱包模式和当前状态执行补偿。
func (s *WalletService) RecoverOrder(ctx context.Context, order *model.GameOrder) error {
	if order == nil {
		return pkg.ErrOrderNotFound
	}
	switch order.Status {
	case model.OrderStatusPending:
		return s.recoverPendingRollback(ctx, order)
	case model.OrderStatusBetSuccess:
		return s.recoverGameResult(ctx, order)
	case model.OrderStatusWaitSettle:
		return s.recoverSettle(ctx, order)
	case model.OrderStatusWaitRollback:
		return s.recoverRollback(ctx, order)
	default:
		return nil
	}
}

// recoverPendingRollback 通过 Operator 的幂等取消语义消除下注结果不确定性。
func (s *WalletService) recoverPendingRollback(ctx context.Context, order *model.GameOrder) error {
	if order.WalletMode != model.WalletModeSingle {
		return pkg.ErrWalletModeInvalid
	}
	agent, user, game, err := s.recoveryEntities(order)
	if err != nil {
		return err
	}
	resp, err := s.operatorClient.Rollback(
		ctx, agent.OperatorURL, agent, user.Username, order.OrderNo,
		game.GameCode, pkg.ToAmount(order.BetAmount),
	)
	if err != nil {
		return err
	}
	if resp == nil {
		return pkg.NewError(pkg.OPERATOR_ROLLBACK_FAILED, "operator rollback returned nil response")
	}
	return s.orderRepo.MarkPendingRolledBack(
		ctx, order.ID, pkg.ToMoney(resp.Balance), resp.Currency,
	)
}

func (s *WalletService) recoverGameResult(ctx context.Context, order *model.GameOrder) error {
	pbReq := &slotpb.SpinReq{
		RequestId:  order.RequestID,
		OrderNo:    order.OrderNo,
		Uid:        order.UID,
		AgentId:    order.AgentID,
		GameId:     order.GameID,
		BetAmount:  order.BetAmount,
		Currency:   order.Currency,
		SpinType:   uint32(order.SpinType),
		FreeSpinId: order.FreeSpinID,
	}

	var (
		resp *slotpb.SpinResp
		err  error
	)
	if order.SpinType == model.SpinTypeFreeSpin {
		resp, err = s.slotAdapter.FreeSpin(ctx, pbReq)
	} else {
		resp, err = s.slotAdapter.Spin(ctx, pbReq)
	}
	if err != nil {
		var rpcErr *skynet.Error
		if !errors.As(err, &rpcErr) {
			return err
		}
		if order.SpinType == model.SpinTypeFreeSpin {
			return s.orderRepo.UpdateStatus(ctx, order.ID, model.OrderStatusBetSuccess, model.OrderStatusFailed)
		}
		if updateErr := s.orderRepo.UpdateWaitRollback(ctx, order.ID, rpcErr.Msg); updateErr != nil {
			return updateErr
		}
		order.Status = model.OrderStatusWaitRollback
		order.RollbackReason = rpcErr.Msg
		return s.recoverRollback(ctx, order)
	}
	if resp == nil || resp.WinAmount < 0 {
		return pkg.NewError(pkg.SPIN_FAILED, "invalid skynet recovery response")
	}
	betAmount := order.BetAmount
	if order.SpinType == model.SpinTypeFreeSpin {
		if resp.FreeSpinId != order.FreeSpinID || resp.BetAmount <= 0 {
			return pkg.NewError(pkg.SPIN_FAILED, "invalid free spin recovery response")
		}
		betAmount = resp.BetAmount
	}
	profit := resp.WinAmount - betAmount
	if order.SpinType == model.SpinTypeFreeSpin {
		profit = resp.WinAmount
	}
	if err := s.orderRepo.UpdateGameResult(
		ctx, order.ID, resp.RoundId, betAmount, resp.WinAmount, profit,
		order.SpinType, resp.FreeSpinId, resp.FreeSpinIndex,
	); err != nil {
		return err
	}
	order.Status = model.OrderStatusWaitSettle
	order.RoundID = resp.RoundId
	order.BetAmount = betAmount
	order.WinAmount = resp.WinAmount
	order.Profit = profit
	order.FreeSpinIndex = resp.FreeSpinIndex
	return s.recoverSettle(ctx, order)
}

func (s *WalletService) recoverSettle(ctx context.Context, order *model.GameOrder) error {
	switch order.WalletMode {
	case model.WalletModeTransfer:
		_, err := s.settleTransferSpin(ctx, order.ID)
		if err == nil {
			s.deleteCache(ctx, order.UID)
		}
		return err

	case model.WalletModeSingle:
		agent, user, game, err := s.recoveryEntities(order)
		if err != nil {
			return err
		}
		resp, err := s.operatorClient.Settle(
			ctx, agent.OperatorURL, agent, user.Username, order.OrderNo,
			game.GameCode, pkg.ToAmount(order.WinAmount),
		)
		if err != nil {
			return err
		}
		if resp == nil {
			return pkg.ErrTransferOrderFailed
		}
		return s.orderRepo.UpdateSettled(
			ctx, order.ID, order.BalanceBefore, pkg.ToMoney(resp.Balance), resp.Currency,
		)

	default:
		return pkg.ErrWalletModeInvalid
	}
}

func (s *WalletService) recoverRollback(ctx context.Context, order *model.GameOrder) error {
	switch order.WalletMode {
	case model.WalletModeTransfer:
		return s.rollbackTransferSpin(ctx, order.ID, order.RollbackReason)

	case model.WalletModeSingle:
		agent, user, game, err := s.recoveryEntities(order)
		if err != nil {
			return err
		}
		resp, err := s.operatorClient.Rollback(
			ctx, agent.OperatorURL, agent, user.Username, order.OrderNo,
			game.GameCode, pkg.ToAmount(order.BetAmount),
		)
		if err != nil {
			return err
		}
		if resp == nil {
			return pkg.ErrTransferOrderFailed
		}
		return s.orderRepo.UpdateRolledBack(
			ctx, order.ID, pkg.ToMoney(resp.Balance), resp.Currency,
		)

	default:
		return pkg.ErrWalletModeInvalid
	}
}

func (s *WalletService) recoveryEntities(order *model.GameOrder) (*model.Agent, *model.User, *model.Game, error) {
	agent, err := s.agentRepo.GetByID(order.AgentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, pkg.ErrAgentNotFound
	}
	if err != nil {
		return nil, nil, nil, err
	}
	user, err := s.userRepo.GetByUID(order.UID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, pkg.ErrUserNotFound
	}
	if err != nil {
		return nil, nil, nil, err
	}
	game, err := s.gameRepo.GetByID(uint64(order.GameID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, pkg.ErrGameNotFound
	}
	if err != nil {
		return nil, nil, nil, err
	}
	return agent, user, game, nil
}
