package mock

import (
	dto "game-api/internal/dto/mock"
	"game-api/internal/service"
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

// 多数 Operator 都是四个接口
type WalletController struct {
	service *service.MockWalletService
}

func NewWalletController(
	service *service.MockWalletService,
) *WalletController {
	return &WalletController{
		service: service,
	}
}

// 查询玩家余额
func (c *WalletController) Balance(ctx *gin.Context) {
	var req dto.BalanceReq

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

// 下注
func (c *WalletController) Bet(ctx *gin.Context) {
	var req dto.BetReq

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

// 结算
func (c *WalletController) Settle(ctx *gin.Context) {
	var req dto.SettleReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}

	resp, err := c.service.Settle(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

	pkg.Success(ctx, resp)
}

// 取消下注
func (c *WalletController) Rollback(ctx *gin.Context) {
	var req dto.RollbackReq

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
