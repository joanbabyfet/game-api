package service

import (
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
func (s *RollbackLogService) Create(log *model.RollbackLog) error {
	return s.repo.Create(log)
}

// Update 更新回滚日志
func (s *RollbackLogService) Update(log *model.RollbackLog) error {
	return s.repo.Update(log)
}

// Delete 删除回滚日志
func (s *RollbackLogService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *RollbackLogService) GetByID(id uint64) (*model.RollbackLog, error) {
	return s.repo.GetByID(id)
}

// GetByRollbackNo 根据回滚单号查询
func (s *RollbackLogService) GetByRollbackNo(rollbackNo string) (*model.RollbackLog, error) {
	return s.repo.GetByRollbackNo(rollbackNo)
}

// GetByOrderNo 根据订单号查询
func (s *RollbackLogService) GetByOrderNo(orderNo string) ([]model.RollbackLog, error) {
	return s.repo.GetByOrderNo(orderNo)
}

// List 回滚日志列表
func (s *RollbackLogService) List(q repository.RollbackLogQuery) ([]model.RollbackLog, error) {
	return s.repo.List(q)
}