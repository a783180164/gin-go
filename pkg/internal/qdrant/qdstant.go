// pkg/internal/qtstant/qtstant.go
package qdrant

import (
	"fmt"
	"sync"

	"gin-go/configs"
	"github.com/qdrant/go-client/qdrant"
)

var (
	once    sync.Once
	client  *qdrant.Client
	errInit error
)

// Init 初始化 Qdrant 客户端，建议在 main() 一启动就调用
func Init() error {
	once.Do(func() {
		cfg := configs.Get().QDSTANT

		// 创建 Qdrant 客户端
		cl, err := qdrant.NewClient(&qdrant.Config{
			Host:   cfg.Host,
			Port:   cfg.Port,
			APIKey: cfg.ApiKey,
		})
		if err != nil {
			errInit = fmt.Errorf("failed to init Qdrant client: %w", err)
			return
		}

		// 可在此处添加其他客户端选项或健康检查
		// e.g., cl.Health()

		// 赋值全局 client 实例
		client = cl
	})

	return errInit
}

// Instance 返回全局 *qdrant.Client
// 请确保在调用 Init 之前已经运行过 Init()
func Instance() *qdrant.Client {
	if client == nil {
		panic("Qdrant client not initialized: call qtstant.Init() first")
	}
	return client
}
