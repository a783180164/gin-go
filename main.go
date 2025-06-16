// main.go
package main

import (
	"context"
	"fmt"
	"gin-go/pkg/bootstrap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	app, err := bootstrap.Initialize()
	if err != nil {
		panic(fmt.Sprintf("failed to bootstrap app: %v", err))
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Engine,
	}

	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("listen error: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	app.Logger.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Fatal("Server Shutdown:", err)
	}
	app.Logger.Println("Server exiting")
}
