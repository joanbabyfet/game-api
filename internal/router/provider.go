package router

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/adapter"
	"game-api/internal/bootstrap"
	"game-api/internal/client/operator"
	provider "game-api/internal/controller/provider"
	"game-api/internal/repository"
	"game-api/internal/service"
)

func RegisterProvider(r *gin.Engine) {

    api := r.Group("/provider")
	
	// Operator Client
	operatorClient := operator.New()

	// Skynet Client
	skynetClient := bootstrap.SkynetClient

	// Adapter
	userAdapter := adapter.NewUserAdapter(skynetClient)
	walletAdapter := adapter.NewWalletAdapter(skynetClient)
	slotAdapter := adapter.NewSlotAdapter(skynetClient)
	systemAdapter := adapter.NewSystemAdapter(skynetClient)

	// Repository
	agentRepo := repository.NewAgentRepository(bootstrap.DB)
	agentGameRepo := repository.NewAgentGameRepository(bootstrap.DB)
	gameRepo := repository.NewGameRepository(bootstrap.DB)
	gameOrderRepo := repository.NewGameOrderRepository(bootstrap.DB)
	userRepo := repository.NewUserRepository(bootstrap.DB)
	walletRepo := repository.NewWalletRepository(bootstrap.DB)

	// Service
	authService := service.NewAuthService(agentRepo)
	gameService := service.NewGameService(
		gameRepo,
		agentRepo,
		userRepo,
		agentGameRepo,
		slotAdapter,
		authService,
	)
	gameOrderService := service.NewGameOrderService(
		gameOrderRepo,
		gameRepo,
		authService,
	)
	walletService := service.NewWalletService(
		walletRepo,
		agentRepo,
		gameRepo,
		userRepo,
		walletAdapter,
		slotAdapter,
		authService,
		operatorClient,
	)
	userService := service.NewUserService(
		userRepo,
		walletRepo,
		gameRepo,
		agentGameRepo,
		userRepo,
		userAdapter,
		authService,
	)
	systemService := service.NewSystemService(
		authService,
		systemAdapter,
	)
	debugService := service.NewDebugService(
		agentRepo,
	)

	// Controller
	gameController := provider.NewGameController(
		gameService,
	)
	orderController := provider.NewOrderController(
		gameOrderService,
	)
	walletController := provider.NewWalletController(
		walletService,
	)
	userController := provider.NewUserController(
		userService,
	)
	systemController := provider.NewSystemController(
		systemService,
	)
	debugController := provider.NewDebugController(
		debugService,
	)
	
	//单一钱包接口
	api.POST("/player/login", userController.Login)
	api.POST("/game_url", gameController.GetGameURL)
	api.POST("/game_list", gameController.List)
	api.POST("/player/kick", userController.Kick)
	api.POST("/get_order_log", orderController.GetOrderLog)
	api.POST("/ping", systemController.Ping)

	//Cocos 
	api.POST("/player/balance", walletController.Balance) //进入游戏、需要刷新余额时
	api.POST("/spin", walletController.Spin) //每点击一次 Spin
	
	//测试用
	api.POST("/debug/sign", debugController.Sign)
}