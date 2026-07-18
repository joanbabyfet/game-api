package service

import (
	"context"
	"errors"
	"game-api/internal/adapter"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"

	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
	walletRepo *repository.WalletRepository
	gameRepo *repository.GameRepository
	agentGameRepo *repository.AgentGameRepository
	userRepo *repository.UserRepository
	userAdapter *adapter.UserAdapter
	authService *AuthService
}

func NewUserService(
	repo *repository.UserRepository,
	walletRepo *repository.WalletRepository,
	gameRepo *repository.GameRepository,
	agentGameRepo *repository.AgentGameRepository,
	userRepo *repository.UserRepository,
	userAdapter *adapter.UserAdapter,
	authService *AuthService,
) *UserService {
	return &UserService{
		repo: repo,
		walletRepo: walletRepo,
		gameRepo: gameRepo,
		agentGameRepo: agentGameRepo,
		userRepo: userRepo,
		userAdapter: userAdapter,
		authService: authService,
	}
}

// Create 新增玩家
func (s *UserService) Create(user *model.User) error {
	return s.repo.Create(user)
}

// Update 更新玩家
func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

// Delete 删除玩家
func (s *UserService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// GetByUID 根据UID查询
func (s *UserService) GetByUID(uid uint64) (*model.User, error) {
	return s.repo.GetByUID(uid)
}

// List 玩家列表
func (s *UserService) List(q repository.UserQuery) ([]model.User, error) {
	return s.repo.List(q)
}

// 踢某玩家 (skynet)
func (s *UserService) Kick(ctx context.Context, req *provider.KickReq) (*provider.KickResp, error) {

	// 驗證 Agent 簽名
	data := pkg.BuildSignData(map[string]any{
		"app_id":    req.AppID,
		"game_id":   req.GameCode,
		"player_id": req.PlayerID,
		"timestamp": req.Timestamp,
	})
	
	agent, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}

	// 檢查代理是否開通該遊戲
	game, err := s.gameRepo.GetByCode(req.GameCode)
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

	// TODO:
	// 後續依 game_id 路由不同 Game Server
	// 目前只有 Slot，直接呼叫即可

	_, err = s.userAdapter.Kick(ctx, user.UID)
	if err != nil {
		return nil, err
	}

	return &provider.KickResp{}, nil
}

//登录创角
func (s *UserService) Login(ctx context.Context, req *provider.LoginReq) (*provider.LoginResp, error) {

	// 驗證 Agent 簽名
	data := pkg.BuildSignData(map[string]any{
		"app_id":    	req.AppID,
		"game_code": 	req.GameCode,
		"player_id":	req.PlayerID,
		"timestamp": 	req.Timestamp,
	})

	agent, err := s.authService.VerifySign(ctx, req.AppID, data, req.Sign)
	if err != nil {
		return nil, err
	}

	// 檢查代理是否開通該遊戲
	game, err := s.gameRepo.GetByCode(req.GameCode)
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

	now := pkg.Timestamp()

	// 查詢玩家
	user, err := s.repo.GetByAgentAndUsername(agent.ID, req.PlayerID)
	if err != nil {
		// 只有不存在才建立
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// nickname := req.Nickname
		// if nickname == "" {
		// 	nickname = req.PlayerID
		// }

		// 建立玩家
		user = &model.User{
			AgentID:       agent.ID,
			Username:      req.PlayerID,
			Nickname:      req.PlayerID,
			Currency:      agent.Currency, //第1次创建玩家时确定，以后不会改变
			Status:        model.UserStatusEnable,
			LastLoginTime: now,
			CreateTime:    now,
			UpdateTime:    now,
		}

		if err := s.repo.Create(user); err != nil {
			return nil, err
		}

		// 建立錢包
		wallet := &model.Wallet{
			UID:           user.UID,
			AgentID:       agent.ID,
			Balance:       0,
			FreezeBalance: 0,
			CreateTime:    now,
			UpdateTime:    now,
		}

		if err := s.walletRepo.Create(wallet); err != nil {
			return nil, err
		}

	} else {

		// 玩家是否啟用
		if user.Status != model.UserStatusEnable {
			return nil, pkg.ErrForbidden
		}

		// 昵称为空则保留原昵称
		// if req.Nickname != "" {
		// 	user.Nickname = req.Nickname
		// }

		user.LastLoginTime = now
		user.UpdateTime = now

		if err := s.repo.Update(user); err != nil {
			return nil, err
		}
	}
	
	// 建立 Skynet Session
	_, err = s.userAdapter.Login(
		ctx,
		user.UID,
	)
	if err != nil {
		return nil, err
	}

	return &provider.LoginResp{}, nil
}