package main

import (
	"game-api/internal/bootstrap"
	"game-api/internal/router"
)

func main() {

    app := bootstrap.New("configs/provider.yaml")

    router.RegisterProvider(app.Engine)

    app.Engine.Run(":8080")

}