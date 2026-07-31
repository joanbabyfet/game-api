package adapter

import (
	"game-api/internal/client/skynet"
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