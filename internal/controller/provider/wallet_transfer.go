package provider

import (
	providerdto "game-api/internal/dto/provider"
	"game-api/internal/service"
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type WalletTransferController struct {
	service *service.WalletTransferService
}

func NewWalletTransferController(service *service.WalletTransferService) *WalletTransferController {
	return &WalletTransferController{service: service}
}

//转入游戏钱包
func (c *WalletTransferController) TransferIn(ctx *gin.Context) {
	var req providerdto.TransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}
	resp, err := c.service.TransferIn(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}
	pkg.Success(ctx, resp)
}

//转出游戏钱包
func (c *WalletTransferController) TransferOut(ctx *gin.Context) {
	var req providerdto.TransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}
	resp, err := c.service.TransferOut(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}
	pkg.Success(ctx, resp)
}

//查询转账订单状态
func (c *WalletTransferController) Status(ctx *gin.Context) {
	var req providerdto.TransferStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
		return
	}
	resp, err := c.service.Status(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}
	pkg.Success(ctx, resp)
}
