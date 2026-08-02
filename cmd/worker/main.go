package main

import (
	"game-api/internal/adapter"
	"game-api/internal/bootstrap"
	"game-api/internal/client/operator"
	"game-api/internal/cron"
	"game-api/internal/repository"
	"game-api/internal/service"
)

func main() {

	bootstrap.New("configs/worker.yaml")

	agentRepo := repository.NewAgentRepository(bootstrap.DB)
	agentGameRepo := repository.NewAgentGameRepository(bootstrap.DB)
	gameRepo := repository.NewGameRepository(bootstrap.DB)
	userRepo := repository.NewUserRepository(bootstrap.DB)
	walletRepo := repository.NewWalletRepository(bootstrap.DB)
	walletLogRepo := repository.NewWalletLogRepository(bootstrap.DB)
	orderRepo := repository.NewGameOrderRepository(bootstrap.DB)
	rollbackLogRepo := repository.NewRollbackLogRepository(bootstrap.DB)

	authService := service.NewAuthService(agentRepo)
	walletService := service.NewWalletService(
		bootstrap.DB,
		bootstrap.GetRedis(),
		walletRepo,
		agentRepo,
		gameRepo,
		agentGameRepo,
		userRepo,
		walletRepo,
		walletLogRepo,
		orderRepo,
		rollbackLogRepo,
		adapter.NewWalletAdapter(bootstrap.SkynetClient),
		adapter.NewSlotAdapter(bootstrap.SkynetClient),
		authService,
		operator.New(),
	)

	cron.Start(cron.NewGameOrderRecoveryJob(orderRepo, walletService))

	select {}

}
