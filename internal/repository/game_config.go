package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type GameConfigRepository struct {
	db *gorm.DB
}

func NewGameConfigRepository(db *gorm.DB) *GameConfigRepository {
	return &GameConfigRepository{
		db: db,
	}
}

type GameConfigQuery struct {
	GameID    uint32
	ConfigKey string

	Page     int
	PageSize int
}

// Create 新增游戏配置
func (r *GameConfigRepository) Create(config *model.GameConfig) error {
	return r.db.Create(config).Error
}

// Update 更新游戏配置
func (r *GameConfigRepository) Update(config *model.GameConfig) error {
	return r.db.
		Model(&model.GameConfig{}).
		Where("id = ?", config.ID).
		Updates(config).
		Error
}

// Delete 删除游戏配置
func (r *GameConfigRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.GameConfig{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *GameConfigRepository) GetByID(id uint64) (*model.GameConfig, error) {

	var config model.GameConfig

	err := r.db.
		Where("id = ?", id).
		First(&config).
		Error

	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetByGameID 根据游戏ID查询
func (r *GameConfigRepository) GetByGameID(gameID uint32) ([]model.GameConfig, error) {

	var configs []model.GameConfig

	err := r.db.
		Where("game_id = ?", gameID).
		Order("id ASC").
		Find(&configs).
		Error

	return configs, err
}

// GetByConfigKey 根据配置Key查询
func (r *GameConfigRepository) GetByConfigKey(
	gameID uint32,
	configKey string,
) (*model.GameConfig, error) {

	var config model.GameConfig

	err := r.db.
		Where("game_id = ?", gameID).
		Where("config_key = ?", configKey).
		First(&config).
		Error

	if err != nil {
		return nil, err
	}

	return &config, nil
}

// List 游戏配置列表
func (r *GameConfigRepository) List(q GameConfigQuery) ([]model.GameConfig, error) {

	var configs []model.GameConfig

	db := r.db

	if q.GameID > 0 {
		db = db.Where("game_id = ?", q.GameID)
	}

	if q.ConfigKey != "" {
		db = db.Where("config_key = ?", q.ConfigKey)
	}

	offset := (q.Page - 1) * q.PageSize

	err := db.
		Order("id DESC").
		Offset(offset).
		Limit(q.PageSize).
		Find(&configs).
		Error

	return configs, err
}