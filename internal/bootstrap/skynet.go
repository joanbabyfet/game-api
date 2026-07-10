package bootstrap

import (
	"fmt"
	"game-api/internal/client/skynet"
	"game-api/internal/config"
)

var SkynetClient *skynet.Client

func InitSkynet() error {

	addr := fmt.Sprintf(
		"%s:%d",
		config.Cfg.Skynet.Host,
		config.Cfg.Skynet.Port,
	)
	client := skynet.New(addr)
	// 建立 TCP 连接
	if err := client.Connect(); err != nil {
		return fmt.Errorf("connect skynet %s failed: %w", addr, err)
	}

	SkynetClient = client

	return nil
}