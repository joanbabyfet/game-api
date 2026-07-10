package service

import (
	"context"
	"game-api/internal/adapter"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
	"game-api/proto/slotpb"
)

type GameService struct {
	repo *repository.GameRepository
	adapter *adapter.SlotAdapter
}

func NewGameService(
	repo *repository.GameRepository,
	adapter *adapter.SlotAdapter,
) *GameService {
	return &GameService{
		repo: repo,
		adapter: adapter,
	}
}

// Create 新增游戏
func (s *GameService) Create(game *model.Game) error {
	return s.repo.Create(game)
}

// Update 更新游戏
func (s *GameService) Update(game *model.Game) error {
	return s.repo.Update(game)
}

// Delete 删除游戏
func (s *GameService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID查询
func (s *GameService) GetByID(id uint64) (*model.Game, error) {
	return s.repo.GetByID(id)
}

// GetByGameID 根据游戏ID查询
func (s *GameService) GetByGameID(gameID uint32) (*model.Game, error) {
	return s.repo.GetByGameID(gameID)
}

// GetByCode 根据游戏编码查询
func (s *GameService) GetByCode(code string) (*model.Game, error) {
	return s.repo.GetByCode(code)
}

// List 游戏列表
func (s *GameService) List(q repository.GameQuery) ([]model.Game, error) {
	//业务参数校验
	if q.Provider == "" {
		return nil, pkg.ErrInvalidParam
	}

	list, err := s.repo.List(q)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 下注
func (s *GameService) Bet(ctx context.Context, req *provider.BetReq) (*provider.BetResp, error) {

	// 解析 JWT
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}

	// 组装 Proto Request
	pbReq := &slotpb.BetReq{
		Uid:     claims.UID,
		AgentId: claims.AgentID,
		GameId:  req.GameID,
		BetAmount:  req.BetAmount,
	}

	// 调用 Skynet
	resp, err := s.adapter.Bet(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	// 转换 Provider DTO
	return &provider.BetResp{
		Balance:  resp.Balance,
	}, nil
}