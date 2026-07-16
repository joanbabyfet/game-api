package service

import (
	"context"

	"game-api/internal/adapter"
	"game-api/internal/dto/provider"
	"game-api/pkg"
)

type SystemService struct {
	authService *AuthService
	adapter     *adapter.SystemAdapter
}

func NewSystemService(
	authService *AuthService,
	adapter *adapter.SystemAdapter,
) *SystemService {
	return &SystemService{
		authService: authService,
		adapter: adapter,
	}
}

//测试 operator 与 provider api 连接是否正常
func (s *SystemService) Ping(ctx context.Context, req *provider.PingReq) (*provider.PingResp, error) {

	data := pkg.BuildSignData(map[string]any{
		"app_id":    req.AppID,
		"timestamp": req.Timestamp,
	})

	_, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}

	_, err = s.adapter.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &provider.PingResp{}, nil
}