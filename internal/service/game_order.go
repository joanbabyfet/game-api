package service

import (
	"game-api/internal/model"
	"game-api/internal/repository"
)

type GameOrderService struct {
	repo *repository.GameOrderRepository
}

func NewGameOrderService(
	repo *repository.GameOrderRepository,
) *GameOrderService {
	return &GameOrderService{
		repo: repo,
	}
}

// Create 新增注单
func (s *GameOrderService) Create(order *model.GameOrder) error {
	return s.repo.Create(order)
}

// Update 更新注单
func (s *GameOrderService) Update(order *model.GameOrder) error {
	return s.repo.Update(order)
}

// Delete 删除注单
func (s *GameOrderService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *GameOrderService) GetByID(id uint64) (*model.GameOrder, error) {
	return s.repo.GetByID(id)
}

// GetByOrderNo 根据订单号查询
func (s *GameOrderService) GetByOrderNo(orderNo string) (*model.GameOrder, error) {
	return s.repo.GetByOrderNo(orderNo)
}

// GetByRequestID 根据RequestID查询
func (s *GameOrderService) GetByRequestID(requestID string) (*model.GameOrder, error) {
	return s.repo.GetByRequestID(requestID)
}

// GetByRoundID 根据RoundID查询
func (s *GameOrderService) GetByRoundID(roundID string) ([]model.GameOrder, error) {
	return s.repo.GetByRoundID(roundID)
}

// List 注单列表
func (s *GameOrderService) List(q repository.GameOrderQuery) ([]model.GameOrder, error) {
	return s.repo.List(q)
}