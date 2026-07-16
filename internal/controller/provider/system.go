package provider

import (
	"game-api/internal/dto/provider"
	"game-api/internal/service"
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type SystemController struct {
	service *service.SystemService
}

func NewSystemController(
	service *service.SystemService,
) *SystemController {
	return &SystemController{
		service: service,
	}
}

//健康检查测试 operator 与 provider api 连接是否正常
func (c *SystemController) Ping(ctx *gin.Context) {

	var req provider.PingReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.Ping(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}