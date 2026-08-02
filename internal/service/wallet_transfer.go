package service

import (
	"context"
	"errors"
	"fmt"
	providerdto "game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type WalletTransferService struct {
	db                *gorm.DB
	redisClient       *redis.Client
	userRepo          *repository.UserRepository
	gameRepo          *repository.GameRepository
	agentGameRepo     *repository.AgentGameRepository
	walletRepo        *repository.WalletRepository
	walletLogRepo     *repository.WalletLogRepository
	transferOrderRepo *repository.WalletTransferRepository
	gameOrderRepo     *repository.GameOrderRepository
	authService       *AuthService
}

func NewWalletTransferService(
	db *gorm.DB,
	redisClient *redis.Client,
	userRepo *repository.UserRepository,
	gameRepo *repository.GameRepository,
	agentGameRepo *repository.AgentGameRepository,
	walletRepo *repository.WalletRepository,
	walletLogRepo *repository.WalletLogRepository,
	transferOrderRepo *repository.WalletTransferRepository,
	gameOrderRepo *repository.GameOrderRepository,
	authService *AuthService,
) *WalletTransferService {
	return &WalletTransferService{
		db: db, redisClient: redisClient, userRepo: userRepo, gameRepo: gameRepo,
		agentGameRepo: agentGameRepo, walletRepo: walletRepo, walletLogRepo: walletLogRepo,
		transferOrderRepo: transferOrderRepo, gameOrderRepo: gameOrderRepo, authService: authService,
	}
}

// 转入游戏钱包
func (s *WalletTransferService) TransferIn(ctx context.Context, req *providerdto.TransferReq) (*providerdto.TransferResp, error) {
	return s.transfer(ctx, req, model.GameTransferTypeIn)
}

// 转出游戏钱包
func (s *WalletTransferService) TransferOut(ctx context.Context, req *providerdto.TransferReq) (*providerdto.TransferResp, error) {
	return s.transfer(ctx, req, model.GameTransferTypeOut)
}

func (s *WalletTransferService) transfer(ctx context.Context, req *providerdto.TransferReq, transferType int8) (*providerdto.TransferResp, error) {
	amount := pkg.ToMoney(req.Amount)
	if amount <= 0 {
		return nil, pkg.ErrInvalidParam
	}

	agent, err := s.authService.VerifySign(ctx, req.AppID, transferSignData(req), req.Sign)
	if err != nil {
		return nil, err
	}
	if agent.WalletMode != model.WalletModeTransfer {
		return nil, pkg.ErrWalletModeInvalid
	}
	if !strings.EqualFold(agent.Currency, req.Currency) {
		return nil, pkg.ErrCurrencyMismatch
	}

	user, err := s.userRepo.GetByAgentAndUsername(agent.ID, req.PlayerID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if user.Status != model.UserStatusEnable {
		return nil, pkg.ErrForbidden
	}

	game, err := s.gameRepo.GetByCode(req.GameCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrGameNotFound
	}
	if err != nil {
		return nil, err
	}
	if game.Status != model.GameStatusEnable {
		return nil, pkg.ErrGameDisabled
	}
	agentGame, err := s.agentGameRepo.GetByAgentGame(agent.ID, game.ID)
	if err != nil || agentGame.Status != model.AgentGameStatusEnable {
		return nil, pkg.ErrForbidden
	}

	if existing, findErr := s.transferOrderRepo.GetByThirdOrderNo(ctx, agent.ID, req.ThirdOrderNo); findErr == nil {
		return replayTransfer(existing, req, user.UID, game.ID, transferType, amount)
	} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, findErr
	}

	now := time.Now().Unix()
	order := &model.WalletTransfer{
		OrderNo: pkg.GenTransferNo(), RequestID: req.RequestID, ThirdOrderNo: req.ThirdOrderNo,
		UID: user.UID, AgentID: agent.ID, GameID: game.ID, TransferType: transferType,
		Amount: amount, Currency: agent.Currency, Status: model.GameTransferStatusPending,
		CreateTime: now, UpdateTime: now,
	}

	var businessErr error
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderRepo := s.transferOrderRepo.WithTx(tx)
		walletRepo := s.walletRepo.WithTx(tx)
		walletLogRepo := s.walletLogRepo.WithTx(tx)
		gameOrderRepo := s.gameOrderRepo.WithTx(tx)

		if err := orderRepo.Create(ctx, order); err != nil {
			return err
		}
		wallet, err := walletRepo.GetByUIDForUpdate(ctx, user.UID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			businessErr = pkg.ErrWalletNotFound
			return orderRepo.UpdateFailed(ctx, order.ID, pkg.WALLET_NOT_FOUND, businessErr.Error())
		}
		if err != nil {
			return err
		}
		if transferType == model.GameTransferTypeOut {
			processing, err := gameOrderRepo.ExistsProcessingTransferOrder(ctx, user.UID, agent.ID)
			if err != nil {
				return err
			}
			if processing {
				businessErr = pkg.ErrOrderProcessing
				return orderRepo.UpdateFailed(ctx, order.ID, pkg.ORDER_PROCESSING, businessErr.Error())
			}
		}

		before := wallet.Balance
		after := before
		logType := model.WalletLogTypeTransferIn
		if transferType == model.GameTransferTypeIn {
			after += amount
			if err := walletRepo.AddBalance(ctx, user.UID, amount); err != nil {
				return err
			}
		} else {
			logType = model.WalletLogTypeTransferOut
			if before < amount {
				businessErr = pkg.ErrInsufficientBalance
				return orderRepo.UpdateFailed(ctx, order.ID, pkg.INSUFFICIENT_BALANCE, businessErr.Error())
			}
			after -= amount
			if err := walletRepo.SubBalance(ctx, user.UID, amount); err != nil {
				return err
			}
		}

		if err := walletLogRepo.Create(ctx, &model.WalletLog{
			UID: user.UID, AgentID: agent.ID, GameID: game.ID, Type: logType,
			Amount: amount, BalanceBefore: before, BalanceAfter: after,
			RefOrderNo: order.OrderNo, CreateTime: now,
		}); err != nil {
			return err
		}
		if err := orderRepo.UpdateSuccess(ctx, order.ID, before, after); err != nil {
			return err
		}
		order.BalanceBefore, order.BalanceAfter = before, after
		order.Status, order.FinishTime, order.UpdateTime = model.GameTransferStatusSuccess, uint32(now), now
		return nil
	})
	if err != nil {
		// A concurrent request may have won the unique-key race.
		if existing, findErr := s.transferOrderRepo.GetByThirdOrderNo(ctx, agent.ID, req.ThirdOrderNo); findErr == nil {
			return replayTransfer(existing, req, user.UID, game.ID, transferType, amount)
		}
		if existing, findErr := s.transferOrderRepo.GetByRequestID(ctx, agent.ID, req.RequestID); findErr == nil {
			return replayTransfer(existing, req, user.UID, game.ID, transferType, amount)
		}
		return nil, err
	}
	if businessErr != nil {
		return nil, businessErr
	}
	s.deleteWalletCache(ctx, user.UID)
	return transferResponse(order), nil
}

func (s *WalletTransferService) Status(ctx context.Context, req *providerdto.TransferStatusReq) (*providerdto.TransferStatusResp, error) {
	data := pkg.BuildSignData(map[string]any{"app_id": req.AppID, "third_order_no": req.ThirdOrderNo, "timestamp": req.Timestamp})
	agent, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}
	if agent.WalletMode != model.WalletModeTransfer {
		return nil, pkg.ErrWalletModeInvalid
	}
	order, err := s.transferOrderRepo.GetByThirdOrderNo(ctx, agent.ID, req.ThirdOrderNo)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pkg.ErrTransferOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	resp := transferResponse(order)
	return resp, nil
}

func transferSignData(req *providerdto.TransferReq) string {
	return pkg.BuildSignData(map[string]any{
		"amount":         req.Amount,
		"app_id":         req.AppID,
		"currency":       req.Currency,
		"game_code":      req.GameCode,
		"player_id":      req.PlayerID,
		"request_id":     req.RequestID,
		"third_order_no": req.ThirdOrderNo,
		"timestamp":      req.Timestamp,
	})
}

func replayTransfer(order *model.WalletTransfer, req *providerdto.TransferReq, uid uint64, gameID uint32, transferType int8, amount int64) (*providerdto.TransferResp, error) {
	if order.RequestID != req.RequestID || order.UID != uid || order.GameID != gameID ||
		order.TransferType != transferType || order.Amount != amount || !strings.EqualFold(order.Currency, req.Currency) {
		return nil, pkg.ErrTransferOrderConflict
	}
	switch order.Status {
	case model.GameTransferStatusSuccess:
		return transferResponse(order), nil
	case model.GameTransferStatusPending:
		return nil, pkg.ErrTransferOrderProcessing
	case model.GameTransferStatusFailed:
		if order.ErrorCode != 0 {
			return nil, pkg.NewError(order.ErrorCode, order.ErrorMessage)
		}
		return nil, pkg.ErrTransferOrderFailed
	default:
		return nil, pkg.ErrTransferOrderConflict
	}
}

func transferResponse(order *model.WalletTransfer) *providerdto.TransferResp {
	return &providerdto.TransferResp{
		OrderNo:      order.OrderNo,
		ThirdOrderNo: order.ThirdOrderNo,
		TransferType: order.TransferType,
		Amount:       pkg.ToAmount(order.Amount),
		//转账后余额
		Balance:  pkg.ToAmount(order.BalanceAfter),
		Currency: order.Currency,
		Status:   order.Status,
	}
}

func (s *WalletTransferService) deleteWalletCache(ctx context.Context, uid uint64) {
	if s.redisClient == nil {
		return
	}
	if err := s.redisClient.Del(fmt.Sprintf("wallet:%d", uid)).Err(); err != nil {
		// Cache invalidation failure must not roll back a committed money transaction.
		return
	}
}
