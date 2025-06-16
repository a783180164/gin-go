package server

import (
	"log"
	"net/http"
)

// RunServer 如果你需要在其他地方启动原生 http 服务，可以使用此函数
func RunServer(addr string, handler http.Handler) {
	log.Printf("Server running at %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
