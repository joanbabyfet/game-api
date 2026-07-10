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

// Authenticate 调用 Skynet 用户认证。
func (a *UserAdapter) Authenticate(
	ctx context.Context,
	uid uint64,
    agentID uint32,
    username string,
) (*userpb.AuthenticateResp, error) {

	pbReq := &userpb.AuthenticateReq{
		Uid:      uid,
		AgentId:  agentID,
		Username: username,
	}

	resp := new(userpb.AuthenticateResp)

	if err := a.client.Call(ctx, skynet.CmdAuthenticate, pbReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Kick 调用 Skynet 踢玩家。
func (a *UserAdapter) Kick(
	ctx context.Context,
	uid uint64,
) (*userpb.KickResp, error) {

	pbReq := &userpb.KickReq{
		Uid: uid,
	}

	resp := new(userpb.KickResp)

	if err := a.client.Call(ctx, skynet.CmdKick, pbReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}