package service

import (
	"context"
	"game-api/internal/adapter"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
	"game-api/proto/walletpb"
)

type WalletService struct {
	repo *repository.WalletRepository
	adapter *adapter.WalletAdapter
}

func NewWalletService(
	repo *repository.WalletRepository,
	adapter *adapter.WalletAdapter,
) *WalletService {
	return &WalletService{
		repo: repo,
		adapter: adapter,
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

// 查询玩家余额
func (s *WalletService) Balance(ctx context.Context, req *provider.BalanceReq) (*provider.BalanceResp, error) {

	// 解析 JWT
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	resp, err := s.adapter.Balance(ctx, claims.UID)
	if err != nil {
		return nil, err
	}

	//做一次转换
	return &provider.BalanceResp{
		Balance: resp.Balance,
		Currency: resp.Currency,
	}, nil
}

func (s *WalletService) Rollback(ctx context.Context, req *provider.RollbackReq) (*provider.RollbackResp, error) {

	// 解析 JWT
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	// 组装 Proto Request
	pbReq := &walletpb.RollbackReq{
		Uid:       claims.UID,
		AgentId:   claims.AgentID,
		OrderNo:   req.OrderNo,
		RequestId: req.RequestID,
		Reason:    req.Reason,
	}

	resp, err := s.adapter.Rollback(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	//做一次转换
	return &provider.RollbackResp{
		Balance: resp.Balance,
	}, nil
}