package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type JackpotRepository struct {
	db *gorm.DB
}

func NewJackpotRepository(db *gorm.DB) *JackpotRepository {
	return &JackpotRepository{
		db: db,
	}
}

type JackpotQuery struct {
	GameID uint32

	Page     int
	PageSize int
}

// Create 新增奖池
func (r *JackpotRepository) Create(pool *model.JackpotPool) error {
	return r.db.Create(pool).Error
}

// Update 更新奖池
func (r *JackpotRepository) Update(pool *model.JackpotPool) error {
	return r.db.
		Model(&model.JackpotPool{}).
		Where("id = ?", pool.ID).
		Updates(pool).
		Error
}

// Delete 删除奖池
func (r *JackpotRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.JackpotPool{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *JackpotRepository) GetByID(id uint64) (*model.JackpotPool, error) {

	var pool model.JackpotPool

	err := r.db.
		Where("id = ?", id).
		First(&pool).
		Error

	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// GetByGameID 根据游戏ID查询
func (r *JackpotRepository) GetByGameID(gameID uint32) (*model.JackpotPool, error) {

	var pool model.JackpotPool

	err := r.db.
		Where("game_id = ?", gameID).
		First(&pool).
		Error

	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// List 奖池列表
func (r *JackpotRepository) List(q JackpotQuery) ([]model.JackpotPool, error) {

	var list []model.JackpotPool

	db := r.db

	if q.GameID > 0 {
		db = db.Where("game_id = ?", q.GameID)
	}

	offset := (q.Page - 1) * q.PageSize

	err := db.
		Order("id DESC").
		Offset(offset).
		Limit(q.PageSize).
		Find(&list).
		Error

	return list, err
}