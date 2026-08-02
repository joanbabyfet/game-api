package service

import (
	"context"
	"errors"
	"game-api/internal/client/skynet"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/pkg"
	"game-api/proto/slotpb"
	"time"

	"gorm.io/gorm"
)

func (s *WalletService) spinTransferWallet(
	ctx context.Context,
	req *provider.SpinReq,
	claims *pkg.JWTClaims,
	agent *model.Agent,
) (*provider.SpinResp, error) {
	game, err := s.gameRepo.GetByCode(req.GameCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrGameNotFound
		}
		return nil, err
	}
	if game.Status != model.GameStatusEnable {
		return nil, pkg.ErrGameDisabled
	}
	if claims.Currency != agent.Currency {
		return nil, pkg.ErrCurrencyMismatch
	}

	isFreeSpin := req.FreeSpinID != ""
	spinType := model.SpinTypeNormal
	if isFreeSpin {
		spinType = model.SpinTypeFreeSpin
	}
	if !isFreeSpin && req.BetAmount <= 0 {
		return nil, pkg.NewError(pkg.INVALID_PARAM, "bet_amount must be greater than zero")
	}
	if isFreeSpin && req.BetAmount != 0 {
		return nil, pkg.NewError(pkg.INVALID_PARAM, "bet_amount must be zero or omitted for free spin")
	}

	if existing, findErr := s.orderRepo.GetByRequestID(ctx, req.RequestID); findErr == nil {
		if err := validateSpinReplay(existing, req, claims, game, spinType, model.WalletModeTransfer); err != nil {
			return nil, err
		}
		return s.replaySpin(existing)
	} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, findErr
	}

	order := &model.GameOrder{
		RequestID:  req.RequestID,
		OrderNo:    pkg.GenOrderNo(),
		UID:        claims.UID,
		AgentID:    agent.ID,
		WalletMode: model.WalletModeTransfer,
		GameID:     game.ID,
		Currency:   agent.Currency,
		SpinType:   spinType,
		FreeSpinID: req.FreeSpinID,
		Status:     model.OrderStatusBetSuccess,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	if !isFreeSpin {
		order.BetAmount = pkg.ToMoney(req.BetAmount)
		if order.BetAmount <= 0 {
			return nil, pkg.ErrInvalidParam
		}
		if err := s.debitTransferSpin(ctx, order); err != nil {
			if existing, findErr := s.orderRepo.GetByRequestID(ctx, req.RequestID); findErr == nil {
				if replayErr := validateSpinReplay(existing, req, claims, game, spinType, model.WalletModeTransfer); replayErr != nil {
					return nil, replayErr
				}
				return s.replaySpin(existing)
			}
			return nil, err
		}
	} else if err := s.orderRepo.Create(ctx, order); err != nil {
		if existing, findErr := s.orderRepo.GetByRequestID(ctx, req.RequestID); findErr == nil {
			if replayErr := validateSpinReplay(existing, req, claims, game, spinType, model.WalletModeTransfer); replayErr != nil {
				return nil, replayErr
			}
			return s.replaySpin(existing)
		}
		return nil, err
	}

	pbReq := &slotpb.SpinReq{
		RequestId:  req.RequestID,
		OrderNo:    order.OrderNo,
		Uid:        claims.UID,
		AgentId:    agent.ID,
		GameId:     game.ID,
		BetAmount:  order.BetAmount,
		Currency:   agent.Currency,
		SpinType:   uint32(spinType),
		FreeSpinId: req.FreeSpinID,
		DebugFail:  req.DebugFail,
	}

	var spinResp *slotpb.SpinResp
	if isFreeSpin {
		spinResp, err = s.slotAdapter.FreeSpin(ctx, pbReq)
	} else {
		spinResp, err = s.slotAdapter.Spin(ctx, pbReq)
	}
	if err != nil {
		var rpcErr *skynet.Error
		if !errors.As(err, &rpcErr) {
			// Result is unknown. Keep BetSuccess for the recovery worker.
			return nil, err
		}
		if isFreeSpin {
			if updateErr := s.orderRepo.UpdateStatus(ctx, order.ID, model.OrderStatusBetSuccess, model.OrderStatusFailed); updateErr != nil {
				return nil, pkg.ErrOrderUpdateFailed
			}
			return nil, err
		}
		if rollbackErr := s.rollbackTransferSpin(ctx, order.ID, rpcErr.Msg); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	if spinResp == nil || spinResp.WinAmount < 0 {
		return nil, pkg.NewError(pkg.SPIN_FAILED, "invalid skynet spin response")
	}
	if isFreeSpin {
		if spinResp.FreeSpinId != req.FreeSpinID || spinResp.BetAmount <= 0 {
			return nil, pkg.NewError(pkg.SPIN_FAILED, "invalid free spin response")
		}
		order.BetAmount = spinResp.BetAmount
		order.FreeSpinIndex = spinResp.FreeSpinIndex
	}

	profit := spinResp.WinAmount - order.BetAmount
	if isFreeSpin {
		profit = spinResp.WinAmount
	}
	if err := s.orderRepo.UpdateGameResult(
		ctx, order.ID, spinResp.RoundId, order.BetAmount, spinResp.WinAmount,
		profit, spinType, spinResp.FreeSpinId, spinResp.FreeSpinIndex,
	); err != nil {
		return nil, pkg.ErrOrderUpdateFailed
	}
	order.RoundID = spinResp.RoundId
	order.WinAmount = spinResp.WinAmount
	order.Profit = profit
	order.Status = model.OrderStatusWaitSettle

	balanceAfter, err := s.settleTransferSpin(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	//Transfer Spin 成功后删除钱包缓存
	s.deleteCache(ctx, claims.UID)

	return &provider.SpinResp{
		Balance:             pkg.ToAmount(balanceAfter),
		Currency:            agent.Currency,
		RoundID:             spinResp.RoundId,
		WinAmount:           pkg.ToAmount(spinResp.WinAmount),
		SpinType:            uint8(spinResp.SpinType),
		FreeSpinID:          spinResp.FreeSpinId,
		FreeSpinIndex:       spinResp.FreeSpinIndex,
		FreeSpinTotalCount:  spinResp.FreeSpinTotalCount,
		FreeSpinRemainCount: spinResp.FreeSpinRemainCount,
	}, nil
}

// 游戏钱包本地下注扣钱 (本地扣款、钱包流水、注单状态在同一事务内)
func (s *WalletService) debitTransferSpin(ctx context.Context, order *model.GameOrder) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		walletRepo := s.walletRepo.WithTx(tx)
		orderRepo := s.orderRepo.WithTx(tx)
		walletLogRepo := s.walletLogRepo.WithTx(tx)

		wallet, err := walletRepo.GetByUIDForUpdate(ctx, order.UID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.ErrWalletNotFound
		}
		if err != nil {
			return err
		}
		if wallet.AgentID != order.AgentID {
			return pkg.ErrForbidden
		}
		if wallet.Balance < order.BetAmount {
			return pkg.ErrInsufficientBalance
		}

		order.Status = model.OrderStatusPending
		order.BalanceBefore = wallet.Balance
		if err := orderRepo.Create(ctx, order); err != nil {
			return err
		}
		if err := walletRepo.SubBalance(ctx, order.UID, order.BetAmount); err != nil {
			return err
		}
		if err := walletLogRepo.Create(ctx, &model.WalletLog{
			UID: order.UID, AgentID: order.AgentID, GameID: order.GameID,
			Type: model.WalletLogTypeBet, Amount: order.BetAmount,
			BalanceBefore: wallet.Balance, BalanceAfter: wallet.Balance - order.BetAmount,
			RefOrderNo: order.OrderNo, CreateTime: time.Now().Unix(),
		}); err != nil {
			return err
		}
		if err := orderRepo.UpdateBetSuccess(ctx, order.ID, wallet.Balance); err != nil {
			return err
		}
		order.Status = model.OrderStatusBetSuccess
		return nil
	})
	if err == nil {
		s.deleteCache(ctx, order.UID)
	}
	return err
}

// 游戏钱包本地派奖加钱 (Skynet成功后，本地派奖与注单结算在同一事务内)
func (s *WalletService) settleTransferSpin(ctx context.Context, orderID uint64) (int64, error) {
	var balanceAfter int64
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		walletRepo := s.walletRepo.WithTx(tx)
		walletLogRepo := s.walletLogRepo.WithTx(tx)

		order, err := orderRepo.GetByIDForUpdate(ctx, orderID)
		if err != nil {
			return err
		}
		if order.Status == model.OrderStatusSettled {
			balanceAfter = order.BalanceAfter
			return nil
		}
		if order.Status != model.OrderStatusWaitSettle {
			return pkg.ErrOrderStatus
		}
		wallet, err := walletRepo.GetByUIDForUpdate(ctx, order.UID)
		if err != nil {
			return err
		}
		before := wallet.Balance
		balanceAfter = before + order.WinAmount
		if order.WinAmount > 0 {
			if err := walletRepo.AddBalance(ctx, order.UID, order.WinAmount); err != nil {
				return err
			}
			logType := model.WalletLogTypeWin
			if order.SpinType == model.SpinTypeFreeSpin {
				logType = model.WalletLogTypeFreeSpin
			}
			if err := walletLogRepo.Create(ctx, &model.WalletLog{
				UID: order.UID, AgentID: order.AgentID, GameID: order.GameID,
				Type: logType, Amount: order.WinAmount,
				BalanceBefore: before, BalanceAfter: balanceAfter,
				RefOrderNo: order.OrderNo, CreateTime: time.Now().Unix(),
			}); err != nil {
				return err
			}
		}
		balanceBefore := order.BalanceBefore
		if order.SpinType == model.SpinTypeFreeSpin {
			balanceBefore = before
		}
		return orderRepo.UpdateSettled(ctx, order.ID, balanceBefore, balanceAfter, order.Currency)
	})
	return balanceAfter, err
}

// Skynet明确业务失败时退回本地下注
func (s *WalletService) rollbackTransferSpin(ctx context.Context, orderID uint64, reason string) error {
	var uid uint64
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderRepo := s.orderRepo.WithTx(tx)
		walletRepo := s.walletRepo.WithTx(tx)
		walletLogRepo := s.walletLogRepo.WithTx(tx)

		order, err := orderRepo.GetByIDForUpdate(ctx, orderID)
		if err != nil {
			return err
		}
		uid = order.UID
		if order.Status == model.OrderStatusRolledBack {
			return nil
		}
		if order.Status != model.OrderStatusBetSuccess {
			return pkg.ErrOrderStatus
		}
		wallet, err := walletRepo.GetByUIDForUpdate(ctx, order.UID)
		if err != nil {
			return err
		}
		before := wallet.Balance
		after := before + order.BetAmount
		if err := walletRepo.AddBalance(ctx, order.UID, order.BetAmount); err != nil {
			return err
		}
		if err := walletLogRepo.Create(ctx, &model.WalletLog{
			UID: order.UID, AgentID: order.AgentID, GameID: order.GameID,
			Type: model.WalletLogTypeRollback, Amount: order.BetAmount,
			BalanceBefore: before, BalanceAfter: after,
			RefOrderNo: order.OrderNo, CreateTime: time.Now().Unix(),
		}); err != nil {
			return err
		}
		return orderRepo.UpdateRolledBackFromBetSuccess(ctx, order.ID, after, order.Currency, reason)
	})
	if err == nil && uid > 0 {
		s.deleteCache(ctx, uid)
	}
	return err
}

func validateSpinReplay(
	order *model.GameOrder,
	req *provider.SpinReq,
	claims *pkg.JWTClaims,
	game *model.Game,
	spinType uint8,
	walletMode int8,
) error {
	if order.UID != claims.UID || order.AgentID != claims.AgentID || order.GameID != game.ID ||
		order.SpinType != spinType || order.FreeSpinID != req.FreeSpinID || order.WalletMode != walletMode {
		return pkg.NewError(pkg.ORDER_STATUS_ERROR, "request_id conflicts with existing spin")
	}
	if spinType == model.SpinTypeNormal && order.BetAmount != pkg.ToMoney(req.BetAmount) {
		return pkg.NewError(pkg.ORDER_STATUS_ERROR, "request_id conflicts with existing bet amount")
	}
	return nil
}
