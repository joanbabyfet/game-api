package provider

import (
	"game-api/internal/dto/provider"
	"game-api/internal/service"

	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
    service *service.GameOrderService
}

func NewOrderController(
	service *service.GameOrderService,
) *OrderController {
	return &OrderController{
		service: service,
	}
}

//获取注单记录
func (c *OrderController) GetOrderLog(ctx *gin.Context) {

	var req provider.OrderLogReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
        return
    }

    resp, err := c.service.GetOrderLog(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

    pkg.Success(ctx, resp)
}