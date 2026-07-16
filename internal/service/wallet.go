package service

import (
	"context"
	"game-api/internal/adapter"
	"game-api/internal/client/operator"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
	"game-api/proto/slotpb"

	"github.com/google/uuid"
)

type WalletService struct {
	repo *repository.WalletRepository
	agentRepo *repository.AgentRepository
	gameRepo *repository.GameRepository
	adapter *adapter.WalletAdapter
	slotAdapter *adapter.SlotAdapter
	authService *AuthService
	operatorClient *operator.Client
}

func NewWalletService(
	repo *repository.WalletRepository,
	agentRepo *repository.AgentRepository,
	gameRepo *repository.GameRepository,
	adapter *adapter.WalletAdapter,
	slotAdapter *adapter.SlotAdapter,
	authService *AuthService,
	
	operatorClient *operator.Client,
) *WalletService {
	return &WalletService{
		repo: repo,
		agentRepo: agentRepo,
		gameRepo: gameRepo,
		adapter: adapter,
		slotAdapter: slotAdapter,
		authService: authService,
		operatorClient: operatorClient,
	}
}

// Create 新增钱包
func (s *WalletService) Create(wallet *model.Wallet) error {
	return s.repo.Create(wallet)
}

// Update 更新钱包
func (s *WalletService) Update(wallet *model.Wallet) error {
	return s.repo.Update(wallet)
}

// Delete 删除钱包
func (s *WalletService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByUID 根据UID查询
func (s *WalletService) GetByUID(uid uint64) (*model.Wallet, error) {
	return s.repo.GetByUID(uid)
}

// List 钱包列表
func (s *WalletService) List(q repository.WalletQuery) ([]model.Wallet, error) {
	return s.repo.List(q)
}

// 查询玩家余额(单一钱包下，余额来源就是 Operator)
func (s *WalletService) Balance(ctx context.Context, req *provider.BalanceReq) (*provider.BalanceResp, error) {

	//解析 JWT (客户端已经登录过，不需要再验 app_id + sign)
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	//获取代理信息
	agent, err := s.agentRepo.GetByID(claims.AgentID)
	if err != nil {
		return nil, err
	}

	//调用 Operator
	resp, err := s.operatorClient.Balance(ctx, agent.OperatorURL, agent, claims.PlayerID)
	if err != nil {
		return nil, err
	}

	// DTO 轉換
	return &provider.BalanceResp{
		Balance: resp.Balance,
		Currency: "USD",
	}, nil
}

// Spin (每点击一次 Spin 调用)
func (s *WalletService) Spin(ctx context.Context, req *provider.SpinReq) (*provider.SpinResp, error) {

	// 解析 JWT
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	// 查询游戏
	game, err := s.gameRepo.GetByCode(req.GameCode)
	if err != nil {
		return nil, pkg.ErrGameNotFound
	}

	//获取代理信息
	agent, err := s.agentRepo.GetByID(claims.AgentID)
	if err != nil {
		return nil, err
	}

	//这些由入口 Provider API 生成
	requestID := uuid.NewString()
	orderNo := pkg.GenOrderNo()
	roundID := pkg.GenRoundID()

	// 1. Operator 扣款
	_, err = s.operatorClient.Bet(
		ctx,
		agent.OperatorURL,
		agent,
		claims.PlayerID,
		orderNo,
		roundID,
		game.GameCode,
		req.BetAmount,
	)
	if err != nil {
		return nil, err
	}

	// 调用 Skynet
	pbReq := &slotpb.SpinReq{
		RequestId: requestID,
		OrderNo:   orderNo,
		RoundId:   roundID,
		Uid:       claims.UID,
		AgentId:   claims.AgentID,
		GameId:    game.ID,
		BetAmount: pkg.ToMoney(req.BetAmount),
	}

	spinResp, err := s.slotAdapter.Spin(ctx, pbReq)
	if err != nil {
		
		// Skynet 失败，回滚扣款
		_, _ = s.operatorClient.Rollback(
			ctx,
			agent.OperatorURL,
			agent,
			claims.PlayerID,
			orderNo,
			roundID,
			game.GameCode,
			req.BetAmount,
		)

		return nil, err
	}

	// 4. Operator 结算
	settleResp, err := s.operatorClient.Settle(
		ctx,
		agent.OperatorURL,
		agent,
		claims.PlayerID,
		orderNo,
		roundID,
		game.GameCode,
		pkg.ToAmount(spinResp.WinAmount),
	)
	if err != nil {
		return nil, err
	}

	return &provider.SpinResp{
		Balance: settleResp.Balance,
		Currency: "USD", //来自 user 表 currency
	}, nil
}

// func (s *WalletService) Rollback(ctx context.Context, req *provider.RollbackReq) (*provider.RollbackResp, error) {

// 	// 解析 JWT
// 	claims, err := pkg.ParseToken(req.Token)
// 	if err != nil {
// 		return nil, pkg.ErrUnauthorized
// 	}

// 	// 组装 Proto Request
// 	pbReq := &walletpb.RollbackReq{
// 		Uid:       claims.UID,
// 		AgentId:   claims.AgentID,
// 		OrderNo:   req.OrderNo,
// 		RequestId: req.RequestID,
// 		Reason:    req.Reason,
// 	}

// 	resp, err := s.adapter.Rollback(ctx, pbReq)
// 	if err != nil {
// 		return nil, err
// 	}

// 	//做一次转换
// 	return &provider.RollbackResp{
// 		Balance: resp.Balance,
// 	}, nil
// }