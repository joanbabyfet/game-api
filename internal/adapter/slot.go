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

// 下注
func (a *SlotAdapter) Bet(
	ctx context.Context,
	req *slotpb.BetReq,
) (*slotpb.BetResp, error) {

	resp := new(slotpb.BetResp)

	if err := a.client.Call(ctx, skynet.CmdBet, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}