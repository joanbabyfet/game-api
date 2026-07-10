package bootstrap

import (
	"game-api/internal/config"
	"game-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type App struct {
    Engine *gin.Engine
}

//三个入口共用同一套初始化
func New(appConfig string) *App {

    // 公共配置
	if err := config.Load("configs/app.yaml"); err != nil {
		panic(err)
	}
	// 应用配置（provider/admin/worker）
	if err := config.Load(appConfig); err != nil {
		panic(err)
	}

	if err := InitMySQL(); err != nil {
        panic(err)
    }

    if err := InitRedis(); err != nil {
        panic(err)
    }

    if err := InitLogger(); err != nil {
        panic(err)
    }

    // 初始化 Skynet Client
	if err := InitSkynet(); err != nil {
        panic(err)
    }

    r := gin.New()

    // 全局 Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

    return &App{
        Engine: r,
    }
}