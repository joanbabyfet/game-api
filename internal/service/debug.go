package service

import (
	"context"

	"game-api/internal/dto/provider"
	"game-api/internal/repository"
	"game-api/pkg"
)

type DebugService struct {
	agentRepo *repository.AgentRepository
}

func NewDebugService(
	agentRepo *repository.AgentRepository,
) *DebugService {
	return &DebugService{
		agentRepo: agentRepo,
	}
}

func (s *DebugService) GenerateSign(
	ctx context.Context,
	req *provider.DebugSignReq,
) (*provider.DebugSignResp, error) {

	appID, ok := req.Fields["app_id"].(string)
	if !ok || appID == "" {
		return nil, pkg.ErrInvalidParam
	}

	agent, err := s.agentRepo.GetByAppID(appID)
	if err != nil {
		return nil, err
	}

	// 組簽名字串
	data := pkg.BuildSignData(req.Fields)
	
	sign := pkg.GenerateSign(data, agent.SecretKey)

	return &provider.DebugSignResp{
		Data: data,
		Sign: sign,
	}, nil
}