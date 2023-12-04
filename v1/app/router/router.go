package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/logic"
	"go_code/lms/app/model"
	"go_code/lms/app/tools"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	i := r.Group("")
	r.Static("/img", "app/img")
	i.Use(checkUser)
	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	//i.GET("/tui", logic.Logout)
	{
		i.GET("/user/create", logic.GetCreate)
		i.POST("/user/create", logic.CreateUser)
	}
	{
		r.GET("/index", logic.Index)
		r.GET("/books", logic.GetBooks)
	}
	{
		i.POST("/books/add", logic.AddBooks)
		i.POST("/books/update", logic.UpdateBooks)
		i.POST("/books/del", logic.DelBooks)
		i.POST("/books/borrow", logic.BorrowBook)
		i.POST("/books/return", logic.ReturnBook)
		r.GET("/image", logic.GetImg)
	}
	{
		r.GET("/captcha", logic.GetCaptcha)
		r.POST("/captcha/verify", func(context *gin.Context) {
			var param tools.CaptchaData
			if err := context.ShouldBind(&param); err != nil {
				context.JSON(http.StatusOK, tools.ParamErr)
				return
			}
			fmt.Printf("参数为：%+v", param)
			if !tools.CaptchaVerify(param) {
				context.JSON(http.StatusOK, tools.ECode{
					Code:    10008,
					Message: "验证失败",
				})
				return
			}
			context.JSON(http.StatusOK, tools.OK)
		})
	}
	r.GET("/", logic.Logout)

	// 创建 HTTP 服务器
	{
		srv := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}

		// 创建一个信号通道，用于接收退出信号
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		// 启动 HTTP 服务器（在协程中）
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}()
		log.Println("Server started")
		// 等待接收退出信号
		<-quit
		log.Println("Received signal to shut down the server")
		// 创建一个上下文对象，设置超时时间为 5 秒
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 关闭服务器（优雅关闭）
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}

		log.Println("Server gracefully stopped")
	}

}
func checkUser(context *gin.Context) {
	var name string
	var id int64 //TODO 存在一个bug
	values := model.GetSession(context)
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}
	if name == "" || id <= 0 {
		context.JSON(http.StatusOK, tools.TouristLogin)
	}
	context.Next()
}
