package provider

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/dto/provider"
	"game-api/internal/model"
	"game-api/internal/repository"
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

// List 游戏列表
func (c *GameController) List(ctx *gin.Context) {

	var req provider.GameListReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	status := model.GameStatusEnable
	query := repository.GameQuery{
		Provider: req.Provider,
		Status:   &status,
		Page:     1,
		PageSize: 20,
	}
	list, err := c.service.List(query)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	//格式化数据
	resp := make([]provider.GameListResp, 0, len(list))
	for _, g := range list {
		resp = append(resp, provider.GameListResp{
			GameCode: g.GameCode,
			Name:     g.GameName,
			Provider: g.Provider,
			Icon:     g.Icon,
			Status:   g.Status,
		})
	}

	pkg.Success(ctx, resp)
}

// 下注
func (c *GameController) Bet(ctx *gin.Context) {

	var req provider.BetReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.Bet(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}