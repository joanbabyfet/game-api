package admin

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/service"
)

type GameController struct {
	service *service.GameService
}

func NewGameController(
	service *service.GameService,
) *GameController {
	return &GameController{
		service: service,
	}
}

// List 游戏列表
func (c *GameController) List(ctx *gin.Context) {

}