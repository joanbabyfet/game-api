package main

import (
	"game-api/internal/bootstrap"
	"game-api/internal/cron"
)

func main() {

    bootstrap.New("configs/worker.yaml")

    cron.Start()

    select {}

}