package provider

import (
	"game-api/internal/dto/provider"
	"game-api/internal/service"
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
    service *service.WalletService
}

func NewWalletController(service *service.WalletService) *WalletController {
	return &WalletController{
		service: service,
	}
}

//获取玩家余额
func (c *WalletController) Balance(ctx *gin.Context) {

	var req provider.BalanceReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.Balance(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}

//取消下注(注单回滚)
func (c *WalletController) Rollback(ctx *gin.Context) {
	var req provider.RollbackReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.Rollback(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}