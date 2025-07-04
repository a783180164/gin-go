package router

import (
	"gin-go/pkg/internal/api/collection"
	"gin-go/pkg/internal/api/ollamatest"
	user "gin-go/pkg/internal/api/user"
	"gin-go/pkg/internal/api/weather"
	"gin-go/pkg/internal/mysql"
	"gin-go/pkg/internal/qdrant"
	jwtmidd "gin-go/pkg/jwt"
	"gin-go/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(r *gin.Engine) {

	jwtCfg := jwtmidd.DefaultJWTConfig()

	users := user.New(logger.Log, mysql.Instance())

	commonApi := r.Group("/api")
	{
		commonApi.POST("/login", users.Login)

	}

	// api模块
	api := r.Group("/api")
	api.Use(jwtCfg.AuthMiddleware())
	{
		api.PUT("/user", users.RegisterUser)
	}

	weathers := weather.New(logger.Log, mysql.Instance())

	// 天气模块
	w := r.Group("weather")
	{
		w.GET("/now", weathers.Now)
	}

	ollama := ollamatest.New(logger.Log, mysql.Instance(), qdrant.Instance())

	o := r.Group("ollamatest")
	{
		o.POST("upload", ollama.Upload)
		o.POST("embed", ollama.EmbedWithOllama)
		o.POST("prompt", ollama.Prompt)
		o.GET("list", ollama.List)
	}

	collect := collection.New(logger.Log, mysql.Instance(), qdrant.Instance())

	c := r.Group("collection")
	{
		c.POST("create", collect.Create)
		c.POST("delete", collect.Delete)
		c.POST("info", collect.Update)
		c.GET("list", collect.List)
	}
}

// pingHandler 一个简单的健康检查
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
