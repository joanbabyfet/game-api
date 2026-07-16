package service

import (
	"context"
	"errors"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"

	"gorm.io/gorm"
)

type AuthService struct {
	agentRepo *repository.AgentRepository
}

func NewAuthService(agentRepo *repository.AgentRepository) *AuthService {
	return &AuthService{agentRepo: agentRepo}
}

//验证签名（AgentCode + Sign）
func (s *AuthService) VerifySign(ctx context.Context, appID string, data string, sign string) (*model.Agent, error) {
	agent, err := s.agentRepo.GetByAppID(appID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrUnauthorized
		}
		return nil, err
	}
	
	//代理是否启用
	if agent.Status != 1 {
		return nil, pkg.ErrForbidden
	}

	// 验证签名
	if !pkg.VerifySign(data, sign, agent.SecretKey) {
		return nil, pkg.ErrUnauthorized
	}
	return agent, nil
}
