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
func (s *AgentService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *AgentService) GetByID(id uint64) (*model.Agent, error) {
	return s.repo.GetByID(id)
}

// GetByAgentID 根据代理ID查询
func (s *AgentService) GetByAgentID(agentID uint32) (*model.Agent, error) {
	return s.repo.GetByAgentID(agentID)
}

// GetByCode 根据代理编码查询
func (s *AgentService) GetByCode(code string) (*model.Agent, error) {
	return s.repo.GetByCode(code)
}

// List 代理列表
func (s *AgentService) List(q repository.AgentQuery) ([]model.Agent, error) {
	return s.repo.List(q)
}