package service

import (
	"context"
	"game-api/internal/model"
	"game-api/internal/repository"
)

type RollbackLogService struct {
	repo *repository.RollbackLogRepository
}

func NewRollbackLogService(
	repo *repository.RollbackLogRepository,
) *RollbackLogService {
	return &RollbackLogService{
		repo: repo,
	}
}

// Create 新增回滚日志
func (s *RollbackLogService) Create(ctx context.Context, log *model.RollbackLog) error {
	return s.repo.Create(ctx, log)
}