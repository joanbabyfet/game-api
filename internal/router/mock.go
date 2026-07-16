package router

import (
	"game-api/internal/bootstrap"
	controller "game-api/internal/controller/mock"
	"game-api/internal/mock"
	"game-api/internal/repository"
	"game-api/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterMock(r *gin.Engine) {

    api := r.Group("/operator")

	// Repository
	agentRepo := repository.NewAgentRepository(bootstrap.DB)

	// Mock 内存钱包
	wallet := mock.NewWallet()

	// Service
	authService := service.NewAuthService(agentRepo)
	mockWalletService := service.NewMockWalletService(
		wallet,
		authService,
	)

	// Controller
	walletController := controller.NewWalletController(mockWalletService)

	//多数 Operator 都是四个接口
	api.POST("/balance", walletController.Balance) //查询玩家余额
	api.POST("/bet", walletController.Bet) //下注
	api.POST("/settle", walletController.Settle) //结算
	api.POST("/rollback", walletController.Rollback) //取消下注
}