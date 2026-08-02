package service

import (
	"game-api/internal/model"
	"game-api/internal/repository"
)

type FreeSpinService struct {
	repo *repository.FreeSpinRepository
}

func NewFreeSpinService(
	repo *repository.FreeSpinRepository,
) *FreeSpinService {
	return &FreeSpinService{
		repo: repo,
	}
}

// Create 新增 Free Spin
func (s *FreeSpinService) Create(freeSpin *model.FreeSpin) error {
	return s.repo.Create(freeSpin)
}

// Update 更新 Free Spin
func (s *FreeSpinService) Update(freeSpin *model.FreeSpin) error {
	return s.repo.Update(freeSpin)
}

// Delete 删除 Free Spin
func (s *FreeSpinService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *FreeSpinService) GetByID(id uint64) (*model.FreeSpin, error) {
	return s.repo.GetByID(id)
}

// GetByFreeSpinID 根据 FreeSpinID 查询
func (s *FreeSpinService) GetByFreeSpinID(freeSpinID string) (*model.FreeSpin, error) {
	return s.repo.GetByFreeSpinID(freeSpinID)
}

// GetByTriggerOrderNo 根据触发注单查询
func (s *FreeSpinService) GetByTriggerOrderNo(orderNo string) (*model.FreeSpin, error) {
	return s.repo.GetByTriggerOrderNo(orderNo)
}

// List Free Spin 列表
func (s *FreeSpinService) List(q repository.FreeSpinQuery) ([]model.FreeSpin, error) {
	return s.repo.List(q)
}