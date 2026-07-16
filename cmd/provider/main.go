package main

import (
	"fmt"
	"game-api/internal/bootstrap"
	"game-api/internal/config"
	"game-api/internal/router"
)

func main() {

    app := bootstrap.New("configs/provider.yaml")

    router.RegisterProvider(app.Engine)

    addr := fmt.Sprintf(":%d", config.Cfg.App.Port)

	if err := app.Engine.Run(addr); err != nil {
		panic(err)
	}

}