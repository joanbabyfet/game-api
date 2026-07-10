package admin

import (
	"github.com/gin-gonic/gin"

	"game-api/internal/dto/admin"
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

	var req admin.GameListReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	query := repository.GameQuery{
		Provider: req.Provider,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	list, err := c.service.List(query)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	//格式化数据
	resp := make([]admin.GameListResp, 0, len(list))
	for _, g := range list {
		resp = append(resp, admin.GameListResp{
			GameCode: g.GameCode,
			Name:     g.GameName,
			Provider: g.Provider,
			Icon:     g.Icon,
			Status:   g.Status,
		})
	}

	pkg.Success(ctx, list)
}