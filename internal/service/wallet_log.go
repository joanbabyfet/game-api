package service

import (
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
func (s *WalletLogService) Create(log *model.WalletLog) error {
	return s.repo.Create(log)
}

// Update 更新钱包流水
func (s *WalletLogService) Update(log *model.WalletLog) error {
	return s.repo.Update(log)
}

// Delete 删除钱包流水
func (s *WalletLogService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *WalletLogService) GetByID(id uint64) (*model.WalletLog, error) {
	return s.repo.GetByID(id)
}

// GetByOrderNo 根据订单号查询
func (s *WalletLogService) GetByOrderNo(orderNo string) ([]model.WalletLog, error) {
	return s.repo.GetByOrderNo(orderNo)
}

// GetByUID 根据UID查询
func (s *WalletLogService) GetByUID(uid uint64) ([]model.WalletLog, error) {
	return s.repo.GetByUID(uid)
}

// List 钱包流水列表
func (s *WalletLogService) List(q repository.WalletLogQuery) ([]model.WalletLog, error) {
	return s.repo.List(q)
}