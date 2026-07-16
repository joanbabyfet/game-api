package adapter

import (
	"context"

	"game-api/internal/client/skynet"
	walletpb "game-api/proto/walletpb"
)

// WalletAdapter 钱包适配器
type WalletAdapter struct {
	client *skynet.Client
}

// NewWalletAdapter 创建钱包适配器
func NewWalletAdapter(client *skynet.Client) *WalletAdapter {
	return &WalletAdapter{
		client: client,
	}
}

// Balance 查询余额
func (a *WalletAdapter) Balance(ctx context.Context, uid uint64) (*walletpb.BalanceResp, error) {

	pbReq := &walletpb.BalanceReq{
		Uid: uid,
	}

	resp := new(walletpb.BalanceResp)

	if err := a.client.Call(ctx, skynet.CmdBalance, pbReq, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Rollback 注单回滚
func (a *WalletAdapter) Rollback(ctx context.Context, req *walletpb.RollbackReq) (*walletpb.RollbackResp, error) {

	resp := new(walletpb.RollbackResp)

	if err := a.client.Call(ctx, skynet.CmdRollback, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}