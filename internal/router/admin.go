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

	//游戏
	gameRepo := repository.NewGameRepository(bootstrap.DB)
	gameService := service.NewGameService(gameRepo, slotAdapter)
	gameController := admin.NewGameController(gameService)

	//测试用
	testController := admin.NewTestController()

    api.GET("/game/list", gameController.List)
	api.GET("/test", testController.Index)
}