package provider

import (
	"game-api/internal/dto/provider"
	"game-api/internal/service"
	"game-api/pkg"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// 踢某玩家 (skynet)
func (c *UserController) Kick(ctx *gin.Context) {

	var req provider.KickReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
        return
    }

    resp, err := c.service.Kick(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

    pkg.Success(ctx, resp)
}

// 登录创角
func (c *UserController) Login(ctx *gin.Context) {

	var req provider.LoginReq

    if err := ctx.ShouldBindJSON(&req); err != nil {
        pkg.Error(ctx, pkg.INVALID_PARAM, err.Error())
        return
    }

    resp, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		pkg.HandleError(ctx, err)
		return
	}

    pkg.Success(ctx, resp)
}