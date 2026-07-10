package provider

import (
	"game-api/internal/dto/provider"
	"game-api/internal/repository"
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

//获取注单历史
func (c *OrderController) History(ctx *gin.Context) {
	var req provider.GameOrderListReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	uid := uint64(10000002)
	agentID := uint64(1)

	query := repository.GameOrderQuery{
		OrderNo:   req.OrderNo,
		RoundID:   req.RoundID,
		UID:       uid,
		AgentID:   agentID,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	list, err := c.service.List(query)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	//格式化数据
	resp := make([]provider.GameOrderListResp, 0, len(list))
	for _, order := range list {
		resp = append(resp, provider.GameOrderListResp{
			OrderNo:    order.OrderNo,
			RoundID:    order.RoundID,
			GameID:     order.GameID,
			BetAmount:  order.BetAmount,
			WinAmount:  order.WinAmount,
			Status:     order.Status,
			CreateTime: order.CreateTime,
		})
	}

	pkg.Success(ctx, resp)
}