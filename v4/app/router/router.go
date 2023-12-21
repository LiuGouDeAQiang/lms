package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
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

func New(cmd *cobra.Command, args []string) {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	i := r.Group("")

	r.Static("/img", "app/img")
	i.Use(checkUser)

	i.Use(logic.Xyz)

	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	//管理员登陆
	r.GET("/admin/Login", logic.GetAdminLogin)
	r.POST("/admin/Login", logic.GetDOAdminLogin)
	//邮箱登录
	r.GET("/emailLogin", logic.GetEmailLogin)
	r.POST("/sending", logic.EmailSend)
	r.POST("/Ok", logic.EmailLogin)
	//手机登录
	r.GET("/telephoneLogin", logic.GetTelephoneLogin)
	r.POST("/sending2", logic.TelephoneSend)
	r.POST("/Ok2", logic.TelephoneLogin)
	//i.GET("/tui", logic.Logout)
	a := r.Group("/admin")
	a.Use(checkAdmin)
	{

		//r.GET("/adminLogin/view", logic.GetAdminView)
		//管理员界面
		{
			a.GET("/view", logic.GetAdminView)
		}
		//管理员列表
		{
			a.GET("/index", onlyAdmin0, logic.GetAdminList0)
			a.GET("/List", onlyAdmin0, logic.GetAdminList)
		}
		//最高级管理员对普通管理员的增删改
		{
			a.POST("/add", onlyAdmin0, logic.AdminAdd)
			a.POST("/del", onlyAdmin0, logic.AdminDel)
			a.POST("/update", onlyAdmin0, logic.AdminUpdate)
		}
		//全部借书记录
		{
			a.GET("/index2", logic.GetInfos0)
			a.GET("/infos", logic.GetInfos)
		}
		//对于图书的增删改查
		{
			a.POST("/books", logic.AddBooks)
			a.PUT("/books/:title", logic.UpdateBooks)
			a.DELETE("/books/:title", logic.DelBooks)
		}
	}
	//用户创建，用户的头像上传
	{
		r.GET("/users/create", logic.GetCreate)
		r.POST("/users", logic.CreateUser)
		r.POST("/users/img/upload", logic.GetUserImg)
	}
	{
		r.GET("/index", logic.Index)
		r.GET("/books", logic.GetBooks)
		r.GET("/mongo", logic.GetBooks)
	}
	{
		i.POST("/books/borrow", onlyUser0, logic.BorrowBook)
		i.POST("/books/return", onlyUser0, logic.ReturnBook)
		i.POST("/books/pay", onlyUser0, tools.HandlePayment)
	}
	{
		r.GET("/books/image", logic.GetImg)
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
	r.GET("/logout", logic.Logout)
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
	var id int64     //TODO 存在一个bug
	var roleId int64 //TODO 存在一个bug
	values := model.GetSession(context)
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}
	if v, ok := values["role_id"]; ok {
		roleId = v.(int64)
	}
	if name == "" || id <= 0 || roleId <= 0 {
		context.JSON(http.StatusOK, tools.TouristLogin)
	}
	context.Next()
}
func checkAdmin(context *gin.Context) {
	var name string
	var id int64     //TODO 存在一个bug
	var roleId int64 //TODO 存在一个bug
	values := model.GetSession(context)
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id = v.(int64)
	}
	if v, ok := values["role_id"]; ok {
		roleId = v.(int64)
	}
	if name == "" || id <= 0 || roleId <= 0 {
		context.Redirect(http.StatusFound, "/login?message=请先登录")
		return
	}
	if roleId == 3 {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "你并不是管理员",
			Data:    nil,
		})
		context.Abort()
	}
	context.Next()
}
func onlyAdmin0(context *gin.Context) {
	var roleId int64 //TODO 存在一个bug
	values := model.GetSession(context)
	if v, ok := values["role_id"]; ok {
		roleId = v.(int64)
	}
	if roleId != 1 {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "你并不是最高级管理员",
			Data:    nil,
		})
		context.Abort()
		return
	}
	context.Next()
}
func onlyUser0(context *gin.Context) {
	var roleId int64 //TODO 存在一个bug
	values := model.GetSession(context)
	if v, ok := values["role_id"]; ok {
		roleId = v.(int64)
	}
	if roleId != 3 {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "只有用户可以进行操作",
			Data:    nil,
		})
		context.Abort()
		return
	}
	context.Next()
}
