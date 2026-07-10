package service

import (
	"game-api/internal/model"
	"game-api/internal/repository"
)

type JackpotService struct {
	repo *repository.JackpotRepository
}

func NewJackpotService(
	repo *repository.JackpotRepository,
) *JackpotService {
	return &JackpotService{
		repo: repo,
	}
}

// Create 新增奖池
func (s *JackpotService) Create(pool *model.JackpotPool) error {
	return s.repo.Create(pool)
}

// Update 更新奖池
func (s *JackpotService) Update(pool *model.JackpotPool) error {
	return s.repo.Update(pool)
}

// Delete 删除奖池
func (s *JackpotService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *JackpotService) GetByID(id uint64) (*model.JackpotPool, error) {
	return s.repo.GetByID(id)
}

// GetByGameID 根据游戏ID查询
func (s *JackpotService) GetByGameID(gameID uint32) (*model.JackpotPool, error) {
	return s.repo.GetByGameID(gameID)
}

// List 奖池列表
func (s *JackpotService) List(q repository.JackpotQuery) ([]model.JackpotPool, error) {
	return s.repo.List(q)
}