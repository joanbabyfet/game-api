package service

import (
	"context"
	"log"

	dto "game-api/internal/dto/mock"
	mockdata "game-api/internal/mock"
	"game-api/pkg"
)

type MockWalletService struct {
	wallet      *mockdata.Wallet
	authService *AuthService
}

func NewMockWalletService(
	wallet *mockdata.Wallet,
	authService *AuthService,
) *MockWalletService {
	return &MockWalletService{
		wallet:      wallet,
		authService: authService,
	}
}

const DefaultCurrency = "USD" //这里测试用, 单一钱包, 余额及币种都是运营商管

// 查询余额
func (s *MockWalletService) Balance(ctx context.Context, req *dto.BalanceReq) (*dto.BalanceResp, error) {

	// 验证签名
	data := pkg.BuildSignData(map[string]any{
		"app_id":    req.AppID,
		"timestamp": req.Timestamp,
		"player_id": req.PlayerID,
	})

	_, err := s.authService.VerifySign(
		ctx,
		req.AppID,
		data,
		req.Sign,
	)
	if err != nil {
		return nil, err
	}

	balance := s.wallet.Balance(req.PlayerID)

	return &dto.BalanceResp{
		Balance:  balance,
		Currency: DefaultCurrency, //在单一钱包模式下, Currency 属于 Operator 的玩家资料，与 Provider API 基本无关
	}, nil
}

// 下注
func (s *MockWalletService) Bet(ctx context.Context, req *dto.BetReq) (*dto.BetResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":     req.AppID,
		"timestamp":  req.Timestamp,
		"player_id":  req.PlayerID,
		"order_no":   req.OrderNo,
		"game_code":  req.GameCode,
		"bet_amount": req.BetAmount,
	})

	log.Println("========== Bet Verify ==========")
	log.Printf("AppID     : %s", req.AppID)
	log.Printf("Sign      : %s", req.Sign)
	log.Printf("Sign Data : %s", data)

	_, err := s.authService.VerifySign(
		ctx,
		req.AppID,
		data,
		req.Sign,
	)
	if err != nil {
		log.Println("========== Verify Failed ==========")
		log.Printf("VerifySign Error: %v", err)
		return nil, err
	}

	balance, err := s.wallet.Bet(req.AppID, req.PlayerID, req.OrderNo, req.GameCode, req.BetAmount)
	if err != nil {
		log.Printf("Wallet Bet Error: %v", err)
		return nil, err
	}

	log.Printf("Wallet Balance: %.2f", balance)

	return &dto.BetResp{
		Balance:  balance,
		Currency: DefaultCurrency,
	}, nil
}

// 结算
func (s *MockWalletService) Settle(ctx context.Context, req *dto.SettleReq) (*dto.SettleResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":     req.AppID,
		"timestamp":  req.Timestamp,
		"player_id":  req.PlayerID,
		"order_no":   req.OrderNo,
		"game_code":  req.GameCode,
		"win_amount": req.WinAmount,
	})

	_, err := s.authService.VerifySign(
		ctx,
		req.AppID,
		data,
		req.Sign,
	)
	if err != nil {
		return nil, err
	}

	balance, err := s.wallet.Settle(req.AppID, req.PlayerID, req.OrderNo, req.GameCode, req.WinAmount)
	if err != nil {
		return nil, err
	}

	return &dto.SettleResp{
		Balance:  balance,
		Currency: DefaultCurrency,
	}, nil
}

// 取消下注
func (s *MockWalletService) Rollback(ctx context.Context, req *dto.RollbackReq) (*dto.RollbackResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":     req.AppID,
		"timestamp":  req.Timestamp,
		"player_id":  req.PlayerID,
		"order_no":   req.OrderNo,
		"game_code":  req.GameCode,
		"bet_amount": req.BetAmount,
	})

	_, err := s.authService.VerifySign(
		ctx,
		req.AppID,
		data,
		req.Sign,
	)
	if err != nil {
		return nil, err
	}

	balance, err := s.wallet.Rollback(req.AppID, req.PlayerID, req.OrderNo, req.GameCode, req.BetAmount)
	if err != nil {
		return nil, err
	}

	return &dto.RollbackResp{
		Balance:  balance,
		Currency: DefaultCurrency,
	}, nil
}
