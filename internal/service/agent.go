package service

import (
	"game-api/internal/model"
	"game-api/internal/repository"
)

type AgentService struct {
	repo *repository.AgentRepository
}

func NewAgentService(
	repo *repository.AgentRepository,
) *AgentService {
	return &AgentService{
		repo: repo,
	}
}

// Create 新增代理
func (s *AgentService) Create(agent *model.Agent) error {
	return s.repo.Create(agent)
}

// Update 更新代理
func (s *AgentService) Update(agent *model.Agent) error {
	return s.repo.Update(agent)
}

// Delete 删除代理
func (s *AgentService) Delete(id uint32) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *AgentService) GetByID(id uint32) (*model.Agent, error) {
	return s.repo.GetByID(id)
}

// List 代理列表
func (s *AgentService) List(q repository.AgentQuery) ([]model.Agent, error) {
	return s.repo.List(q)
}