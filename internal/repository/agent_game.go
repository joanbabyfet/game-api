package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type AgentGameRepository struct {
	db *gorm.DB
}

func NewAgentGameRepository(db *gorm.DB) *AgentGameRepository {
	return &AgentGameRepository{
		db: db,
	}
}

type AgentGameQuery struct {
	AgentID uint32
	GameID  uint32
	Status  *int8

	Page     int
	PageSize int
}

// 获取列表
func (r *AgentGameRepository) List(q *AgentGameQuery) ([]model.AgentGame, error) {

	db := r.db.Model(&model.AgentGame{})

	if q.AgentID > 0 {
		db = db.Where("agent_id = ?", q.AgentID)
	}

	if q.GameID > 0 {
		db = db.Where("game_id = ?", q.GameID)
	}

	if q.Status != nil {
		db = db.Where("status = ?", q.Status)
	}

	db = db.Order("game_id ASC")

	if q.Page > 0 && q.PageSize > 0 {
		offset := (q.Page - 1) * q.PageSize
		db = db.Offset(offset).Limit(q.PageSize)
	}

	var list []model.AgentGame
	err := db.Find(&list).Error

	return list, err
}

// Create
func (r *AgentGameRepository) Create(agentGame *model.AgentGame) error {
	return r.db.Create(agentGame).Error
}

// Update
func (r *AgentGameRepository) Update(agentGame *model.AgentGame) error {
	return r.db.
		Model(&model.AgentGame{}).
		Where("id = ?", agentGame.ID).
		Updates(agentGame).
		Error
}

// Delete
func (r *AgentGameRepository) Delete(id uint64) error {
	return r.db.
		Delete(&model.AgentGame{}, "id = ?", id).
		Error
}

// GetByID
func (r *AgentGameRepository) GetByID(id uint64) (*model.AgentGame, error) {

	var agentGame model.AgentGame

	err := r.db.
		Where("id = ?", id).
		First(&agentGame).
		Error

	if err != nil {
		return nil, err
	}

	return &agentGame, nil
}

// GetByAgentGame
func (r *AgentGameRepository) GetByAgentGame(agentID uint32, gameID uint32) (*model.AgentGame, error) {

	var agentGame model.AgentGame

	err := r.db.
		Where("agent_id = ? AND game_id = ?", agentID, gameID).
		First(&agentGame).
		Error

	if err != nil {
		return nil, err
	}

	return &agentGame, nil
}