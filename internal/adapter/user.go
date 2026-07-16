package adapter

import (
	"context"

	"game-api/internal/client/skynet"
	userpb "game-api/proto/userpb"
)

// UserAdapter 封装用户相关 RPC。
type UserAdapter struct {
	client *skynet.Client
}

func NewUserAdapter(client *skynet.Client) *UserAdapter {
	return &UserAdapter{
		client: client,
	}
}

// Kick 调用 Skynet 踢玩家。
func (a *UserAdapter) Kick(ctx context.Context, uid uint64) (*userpb.KickResp, error) {

	pbReq := &userpb.KickReq{
		Uid: uid,
	}

	resp := new(userpb.KickResp)

	if err := a.client.Call(ctx, skynet.CmdKick, pbReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Login 调用 Skynet 玩家上线
func (a *UserAdapter) Login(ctx context.Context, uid uint64) (*userpb.LoginResp, error) {

	pbReq := &userpb.LoginReq{
		Uid: uid,
	}

	resp := new(userpb.LoginResp)

	if err := a.client.Call(ctx, skynet.CmdLogin, pbReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}