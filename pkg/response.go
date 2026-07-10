package pkg

import (
	"game-api/internal/client/skynet"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//统一返回
type Response struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Success 成功
func Success(c *gin.Context, data interface{}) {

	if data == nil {
		data = struct{}{}
	}

	c.JSON(http.StatusOK, Response{
		Code:      SUCCESS,
		Msg:       "success",
		Timestamp: time.Now().Unix(),
		Data:      data,
	})
}

// Error 错误
func Error(c *gin.Context, code int, msg string) {

	if code == 0 {
		code = UNKNOWN_ERROR
	}

	if msg == "" {
		msg = "error"
	}

	c.JSON(http.StatusOK, Response{
		Code:      code,
		Msg:       msg,
		Timestamp: time.Now().Unix(),
		Data:      struct{}{},
	})
}

// HandleError 根据 error 返回统一错误响应
func HandleError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	// HTTP/App 错误
	if appErr, ok := err.(*AppError); ok {
		Error(ctx, appErr.Code, appErr.Msg)
		return
	}

	// RPC(Skynet) 错误
	if rpcErr, ok := err.(*skynet.Error); ok {
		Error(ctx, int(rpcErr.Code), rpcErr.Msg)
		return
	}

	Error(ctx, UNKNOWN_ERROR, "internal server error")
}