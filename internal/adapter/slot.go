package adapter

import (
	"context"

	"game-api/internal/client/skynet"
	"game-api/proto/slotpb"
)

// SlotAdapter slot适配器
type SlotAdapter struct {
	client *skynet.Client
}

// NewSlotAdapter 创建slot适配器
func NewSlotAdapter(client *skynet.Client) *SlotAdapter {
	return &SlotAdapter{
		client: client,
	}
}

// 旋通旋转
func (a *SlotAdapter) Spin(ctx context.Context, req *slotpb.SpinReq) (*slotpb.SpinResp, error) {

	resp := new(slotpb.SpinResp)

	if err := a.client.Call(ctx, skynet.CmdSpin, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}