package main

import (
	"game-api/internal/bootstrap"
	"game-api/internal/router"
)

func main() {

    app := bootstrap.New("configs/admin.yaml")

    router.RegisterAdmin(app.Engine)

    app.Engine.Run(":8081")

}