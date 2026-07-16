package provider

import (
	"game-api/internal/dto/provider"
	"game-api/internal/service"
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type DebugController struct {
	service *service.DebugService
}

func NewDebugController(
	service *service.DebugService,
) *DebugController {
	return &DebugController{
		service: service,
	}
}

func (c *DebugController) Sign(ctx *gin.Context) {

	var req provider.DebugSignReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.GenerateSign(
		ctx.Request.Context(),
		&req,
	)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}