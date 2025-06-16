// pkg/bootstrap/bootstrap.go
package bootstrap

import (
	"fmt"
	"gin-go/configs"
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/qdrant"
	"gin-go/pkg/logger"
	"gin-go/pkg/router"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Engine *gin.Engine
	Logger *logrus.Logger
}

// Initialize 初始化所有组件：配置、数据库、日志、路由等
func Initialize() (*App, error) {
	// 初始化配置
	config := configs.Get()

	// 初始化日志系统
	log, err := logger.InitDateLogger(
		"logs/app-%Y-%m-%d.log",
		"logs/app.log",
		7,
		logrus.InfoLevel,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	// 初始化数据库
	if err := mysql.Init(); err != nil {
		log.Errorf("failed to init mysql: %v", err)
		return nil, err
	}

	// 初始化
	if err := qdrant.Init(); err != nil {
		log.Errorf("failed to init qdrant: %s", err)
	}

	// 创建 Gin 实例
	r := gin.New()
	enableRequestLog := config.Logger.EnableRequestLog || false
	enableResponseLog := config.Logger.EnableResponseLog || false
	r.Use(logger.GinLogger(log, enableRequestLog, enableResponseLog), gin.Recovery())

	// 注册路由
	// router.RegisterRoutes(r)

	router.SetApiRouter(r)

	return &App{
		Engine: r,
		Logger: log,
	}, nil
}
