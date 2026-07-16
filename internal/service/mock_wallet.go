package service

import (
	"context"

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
		Balance: balance,
		Currency: "USD", //在单一钱包模式下, Currency 属于 Operator 的玩家资料，与 Provider API 基本无关
	}, nil
}

// 下注
func (s *MockWalletService) Bet(ctx context.Context, req *dto.BetReq) (*dto.BetResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":    	req.AppID,
		"timestamp": 	req.Timestamp,
		"player_id":    req.PlayerID,
		"order_no":  	req.OrderNo,
		"round_id":  	req.RoundID,
		"game_code":   	req.GameCode,
		"bet_amount":	req.BetAmount,
		"bet_time":		req.BetTime,
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

	balance, err := s.wallet.Bet(req.PlayerID, req.BetAmount)
	if err != nil {
		return nil, err
	}

	return &dto.BetResp{
		Balance: balance,
		Currency: "USD",
	}, nil
}

// 结算
func (s *MockWalletService) Settle(ctx context.Context, req *dto.SettleReq) (*dto.SettleResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":    	req.AppID,
		"timestamp": 	req.Timestamp,
		"player_id":    req.PlayerID,
		"order_no":  	req.OrderNo,
		"round_id":  	req.RoundID,
		"game_code":   	req.GameCode,
		"bet_amount":	req.BetAmount,
		"win_amount":	req.WinAmount,
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

	balance, err := s.wallet.Settle(req.PlayerID, req.WinAmount)
	if err != nil {
		return nil, err
	}

	return &dto.SettleResp{
		Balance: balance,
		Currency: "USD",
	}, nil
}
// 取消下注
func (s *MockWalletService) Rollback(ctx context.Context, req *dto.RollbackReq) (*dto.RollbackResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":    	req.AppID,
		"timestamp": 	req.Timestamp,
		"player_id":    req.PlayerID,
		"order_no":  	req.OrderNo,
		"round_id":  	req.RoundID,
		"game_code":   	req.GameCode,
		"bet_amount":	req.BetAmount,
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

	balance, err := s.wallet.Rollback(req.PlayerID, req.BetAmount)
	if err != nil {
		return nil, err
	}

	return &dto.RollbackResp{
		Balance: balance,
		Currency: "USD",
	}, nil
}