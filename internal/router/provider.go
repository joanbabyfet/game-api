package router

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/adapter"
	"game-api/internal/bootstrap"
	provider "game-api/internal/controller/provider"
	"game-api/internal/repository"
	"game-api/internal/service"
)

func RegisterProvider(r *gin.Engine) {

    api := r.Group("/provider")
	
	// Skynet Client
	skynetClient := bootstrap.SkynetClient

	// Adapter
	userAdapter := adapter.NewUserAdapter(skynetClient)
	walletAdapter := adapter.NewWalletAdapter(skynetClient)
	slotAdapter := adapter.NewSlotAdapter(skynetClient)

	//游戏
	gameRepo := repository.NewGameRepository(bootstrap.DB)
	gameService := service.NewGameService(gameRepo, slotAdapter)
	gameController := provider.NewGameController(gameService)

	// 历史注单
	gameOrderRepo := repository.NewGameOrderRepository(bootstrap.DB)
	gameOrderService := service.NewGameOrderService(gameOrderRepo)
	orderController := provider.NewOrderController(gameOrderService)

	//钱包
	walletRepo := repository.NewWalletRepository(bootstrap.DB)
	walletService := service.NewWalletService(walletRepo, walletAdapter)
	walletController := provider.NewWalletController(walletService)

	// 玩家
	userRepo := repository.NewUserRepository(bootstrap.DB)
	userService := service.NewUserService(userRepo, walletRepo, userAdapter)
	userController := provider.NewUserController(userService)

    api.GET("/game/list", gameController.List)
	api.GET("/history", orderController.History)
	api.POST("/authenticate", userController.Authenticate)
	api.POST("/player/kick", userController.Kick)
	api.POST("/balance", walletController.Balance)
	api.POST("/bet", gameController.Bet)
	api.POST("/rollback", walletController.Rollback)
}