package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/logic"
	"go_code/lms/app/model"
	"go_code/lms/app/tools"
	"net/http"
)

func New() {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	i := r.Group("")
	i.Static("/img", "app/img")

	i.Use(checkUser)

	r.GET("/login", logic.GetLogin)
	r.POST("/login", logic.DoLogin)
	r.GET("/logout", logic.Logout)
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

	r.Run(":8080")
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
