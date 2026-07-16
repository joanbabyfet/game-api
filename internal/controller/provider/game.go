package provider

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/dto/provider"
	"game-api/internal/service"
	"game-api/pkg"
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

// 获取游戏列表
func (c *GameController) List(ctx *gin.Context) {

	var req provider.GameListReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.List(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}

// 获取进入游戏地址
func (c *GameController) GetGameURL(ctx *gin.Context) {

	var req provider.GameURLReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.GetGameURL(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}