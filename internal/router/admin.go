package router

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/adapter"
	"game-api/internal/bootstrap"
	admin "game-api/internal/controller/admin"
	"game-api/internal/repository"
	"game-api/internal/service"
)

func RegisterAdmin(r *gin.Engine) {

    api := r.Group("/admin")

	// Skynet Client
	skynetClient := bootstrap.SkynetClient

	// Adapter
	slotAdapter := adapter.NewSlotAdapter(skynetClient)

	// Repository
	agentRepo := repository.NewAgentRepository(bootstrap.DB)
	gameRepo := repository.NewGameRepository(bootstrap.DB)
	userRepo := repository.NewUserRepository(bootstrap.DB)
	agentGameRepo := repository.NewAgentGameRepository(bootstrap.DB)

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

	// Controller
	gameController := admin.NewGameController(
		gameService,
	)
	testController := admin.NewTestController()

    api.GET("/game/list", gameController.List)
	api.GET("/test", testController.Index)
}