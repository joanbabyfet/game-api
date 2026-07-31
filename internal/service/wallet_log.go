package service

import (
	"context"
	"game-api/internal/model"
	"game-api/internal/repository"
)

type WalletLogService struct {
	repo *repository.WalletLogRepository
}

func NewWalletLogService(
	repo *repository.WalletLogRepository,
) *WalletLogService {
	return &WalletLogService{
		repo: repo,
	}
}

// Create 新增钱包流水
func (s *WalletLogService) Create(ctx context.Context, log *model.WalletLog) error {
	return s.repo.Create(ctx, log)
}