package tools

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefulShutdown 完成优雅退出
//
//	func GracefulShutdown(server *http.Server) {
//		// 创建一个用于接收终止信号的通道
//		quit := make(chan os.Signal, 1)
//		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//
//		// 等待终止信号
//		<-quit
//
//		// 创建一个上下文对象，设置超时时间
//		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//		defer cancel()
//
//		// 关闭服务器，并等待所有连接关闭
//		if err := server.Shutdown(ctx); err != nil {
//			log.Fatalf("Server shutdown error: %v", err)
//		}
//		server.Close()
//		log.Println("Server gracefully stopped")
//	}
func GracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	// 显示地停止服务器
	if err := server.Close(); err != nil {
		log.Fatalf("Server close error: %v", err)
	}

	log.Println("Server gracefully stopped")
}
