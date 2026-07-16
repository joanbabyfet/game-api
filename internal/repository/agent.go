package repository

import (
	"game-api/internal/model"

	"gorm.io/gorm"
)

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{
		db: db,
	}
}

type AgentQuery struct {
	AgentID uint64
	Code    string
	Name    string
	Status  int8

	Page     int
	PageSize int
}

// Create 新增代理
func (r *AgentRepository) Create(agent *model.Agent) error {
	return r.db.Create(agent).Error
}

// Update 更新代理
func (r *AgentRepository) Update(agent *model.Agent) error {
	return r.db.
		Model(&model.Agent{}).
		Where("id = ?", agent.ID).
		Updates(agent).
		Error
}

// Delete 删除代理
func (r *AgentRepository) Delete(id uint32) error {
	return r.db.
		Delete(&model.Agent{}, "id = ?", id).
		Error
}

// GetByID 根据ID查询
func (r *AgentRepository) GetByID(id uint32) (*model.Agent, error) {

	var agent model.Agent

	err := r.db.
		Where("id = ?", id).
		First(&agent).
		Error

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// GetByAppID 根据 AppID 查询代理
func (r *AgentRepository) GetByAppID(appID string) (*model.Agent, error) {

	var agent model.Agent

	err := r.db.
		Where("app_id = ?", appID).
		First(&agent).
		Error

	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// List 代理列表
func (r *AgentRepository) List(q AgentQuery) ([]model.Agent, error) {

	var agents []model.Agent

	db := r.db

	if q.AgentID > 0 {
		db = db.Where("agent_id = ?", q.AgentID)
	}

	if q.Code != "" {
		db = db.Where("code = ?", q.Code)
	}

	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}

	if q.Status >= 0 {
		db = db.Where("status = ?", q.Status)
	}

	offset := (q.Page - 1) * q.PageSize

	err := db.
		Order("id DESC").
		Offset(offset).
		Limit(q.PageSize).
		Find(&agents).
		Error

	return agents, err
}