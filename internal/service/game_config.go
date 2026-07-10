package service

import (
	"game-api/internal/model"
	"game-api/internal/repository"
)

type GameConfigService struct {
	repo *repository.GameConfigRepository
}

func NewGameConfigService(
	repo *repository.GameConfigRepository,
) *GameConfigService {
	return &GameConfigService{
		repo: repo,
	}
}

// Create 新增游戏配置
func (s *GameConfigService) Create(config *model.GameConfig) error {
	return s.repo.Create(config)
}

// Update 更新游戏配置
func (s *GameConfigService) Update(config *model.GameConfig) error {
	return s.repo.Update(config)
}

// Delete 删除游戏配置
func (s *GameConfigService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *GameConfigService) GetByID(id uint64) (*model.GameConfig, error) {
	return s.repo.GetByID(id)
}

// GetByGameID 根据游戏ID查询
func (s *GameConfigService) GetByGameID(gameID uint32) ([]model.GameConfig, error) {
	return s.repo.GetByGameID(gameID)
}

// GetByConfigKey 根据配置Key查询
func (s *GameConfigService) GetByConfigKey(
	gameID uint32,
	configKey string,
) (*model.GameConfig, error) {
	return s.repo.GetByConfigKey(gameID, configKey)
}

// List 游戏配置列表
func (s *GameConfigService) List(q repository.GameConfigQuery) ([]model.GameConfig, error) {
	return s.repo.List(q)
}