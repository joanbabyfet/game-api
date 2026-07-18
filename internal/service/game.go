package service

import (
	"context"
	"fmt"
	"game-api/internal/adapter"
	"game-api/internal/config"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
	"net/url"
)

type GameService struct {
	repo *repository.GameRepository
	agentRepo *repository.AgentRepository
	userRepo  *repository.UserRepository
	agentGameRepo  *repository.AgentGameRepository
	adapter *adapter.SlotAdapter
	authService *AuthService
}

func NewGameService(
	repo *repository.GameRepository,
	agentRepo *repository.AgentRepository,
	userRepo *repository.UserRepository,
	agentGameRepo  *repository.AgentGameRepository,
	adapter *adapter.SlotAdapter,
	authService *AuthService,
) *GameService {
	return &GameService{
		repo: repo,
		agentRepo: agentRepo,
		userRepo:  userRepo,
		agentGameRepo:  agentGameRepo,
		adapter: adapter,
		authService: authService,
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
func (s *GameService) List(ctx context.Context, req *provider.GameListReq) ([]provider.GameListResp, error) {

	// 验证签名
	data := pkg.BuildSignData(map[string]any{
		"app_id":    req.AppID,
		"timestamp": req.Timestamp,
	})

	agent, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}

	// 查询 Agent 已开通游戏
	status := model.AgentGameStatusEnable
	agentGames, err := s.agentGameRepo.List(&repository.AgentGameQuery{
		AgentID: agent.ID,
		Status:  &status,
	})
	if err != nil {
		return nil, err
	}

	// 提取 GameID
	gameIDs := make([]uint32, 0, len(agentGames))
	for _, item := range agentGames {
		gameIDs = append(gameIDs, item.GameID)
	}

	// 查询游戏
	games, err := s.repo.ListByIDs(gameIDs)
	if err != nil {
		return nil, err
	}

	// 返回
	resp := make([]provider.GameListResp, 0, len(games))

	for _, game := range games {
		resp = append(resp, provider.GameListResp{
			GameCode: game.GameCode,
			GameName: game.GameName,
			Provider: game.Provider,
			Icon:     pkg.ImageURL(game.Icon),
		})
	}

	return resp, nil
}

//获取进入游戏地址
func (s *GameService) GetGameURL(ctx context.Context, req *provider.GameURLReq) (*provider.GameURLResp, error) {

	// 驗證 Agent 簽名
	data := pkg.BuildSignData(map[string]any{
		"app_id":    	req.AppID,
		"game_code": 	req.GameCode,
		"player_id":	req.PlayerID,
		"lang":			req.Lang,
		"timestamp": 	req.Timestamp,
	})

	agent, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}

	// 檢查代理是否開通該遊戲
	game, err := s.repo.GetByCode(req.GameCode)
	if err != nil {
		return nil, pkg.ErrGameNotFound
	}

	agentGame, err := s.agentGameRepo.GetByAgentGame(agent.ID, game.ID)
	if err != nil {
		return nil, pkg.ErrForbidden
	}
	// 遊戲是否啟用
	if agentGame.Status != 1 {
		return nil, pkg.ErrForbidden
	}

	// 查询玩家
	user, err := s.userRepo.GetByAgentAndUsername(agent.ID, req.PlayerID)
	if err != nil {
		return nil, pkg.ErrUserNotFound
	}

	// 生成 JWT
	token, err := pkg.GenerateToken(
		user.UID,
		agent.ID,
		user.Username,
		user.Currency,
	)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("token", token)
	params.Set("game_code", game.GameCode)
	params.Set("lang", req.Lang) 	// 可先固定，之后放到配置

	gameURL := fmt.Sprintf(
		"%s?%s",
		config.Cfg.Provider.GameURL,
		params.Encode(),
	)

	return &provider.GameURLResp{
		URL:      gameURL,
	}, nil
}