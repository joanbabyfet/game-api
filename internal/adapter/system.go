package adapter

import (
	"context"

	"game-api/internal/client/skynet"
	systempb "game-api/proto/systempb"
)

// SystemAdapter 封装系统 RPC。
type SystemAdapter struct {
	client *skynet.Client
}

// NewSystemAdapter 创建 SystemAdapter。
func NewSystemAdapter(client *skynet.Client) *SystemAdapter {
	return &SystemAdapter{
		client: client,
	}
}

// Ping 测试与 Skynet 的连通性。
func (a *SystemAdapter) Ping(ctx context.Context) (*systempb.PingResp, error) {
	pbReq := &systempb.PingReq{}

	resp := &systempb.PingResp{}

	if err := a.client.Call(ctx, skynet.CmdPing, pbReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}