package service

import (
	"context"
	"game-api/internal/adapter"
	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
	"game-api/pkg"
)

type UserService struct {
	repo *repository.UserRepository
	walletRepo *repository.WalletRepository
	userAdapter *adapter.UserAdapter
}

func NewUserService(
	repo *repository.UserRepository,
	walletRepo *repository.WalletRepository,
	userAdapter *adapter.UserAdapter,
) *UserService {
	return &UserService{
		repo: repo,
		walletRepo: walletRepo,
		userAdapter: userAdapter,
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

// GetByUsername 根据用户名查询
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	return s.repo.GetByUsername(username)
}

// List 玩家列表
func (s *UserService) List(q repository.UserQuery) ([]model.User, error) {
	return s.repo.List(q)
}

// Authenticate 登录认证
func (s *UserService) Authenticate(ctx context.Context, req *provider.AuthReq) (*provider.AuthResp, error) {

	// 解析JWT
	claims, err := pkg.ParseToken(req.Token)
	if err != nil {
		return nil, pkg.ErrUnauthorized
	}
	
	resp, err := s.userAdapter.Authenticate(ctx, claims.UID, claims.AgentID, claims.Username)
	if err != nil {
		return nil, err
	}

	//做一次转换
	return &provider.AuthResp{
		UID:      resp.Uid,
		AgentID:  resp.AgentId,
		Username: resp.Username,
		Balance:  resp.Balance,
		Currency: resp.Currency,
	}, nil
}

// 踢某玩家 (skynet)
func (s *UserService) Kick(ctx context.Context, req *provider.KickReq) error {

	_, err := s.userAdapter.Kick(ctx, req.UID)
	if err != nil {
		return err
	}

	return nil
}