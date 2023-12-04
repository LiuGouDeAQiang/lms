package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/model"
	"go_code/lms/app/tools"
	"net/http"
	"os"
	"strconv"
	"time"
)

func AddBooks(context *gin.Context) {
	Title := context.Query("title")
	Num, _ := context.GetPostForm("Num")

	num, err := strconv.ParseInt(Num, 10, 32)
	if err != nil {
		panic(err)
	}
	//构建结构体
	NewBook := model.Books{
		Title: Title,
		Num:   num,

		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}
	if NewBook.Title == " " {
		context.JSON(http.StatusBadRequest, tools.ParamErr)
	}

	if err := model.AddBooks(NewBook); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}
	context.JSON(201, tools.OK)
	return
}

func UpdateBooks(context *gin.Context) {
	Title := context.Query("title")
	Num, _ := context.GetPostForm("Num")
	num, err := strconv.ParseInt(Num, 10, 32)
	if err != nil {
		panic(err)
	}
	//构建结构体
	NewBook := model.Books{
		Title:       Title,
		Num:         num,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}
	if NewBook.Title == " " {
		context.JSON(http.StatusBadRequest, tools.ParamErr)
	}

	if err := model.UpdateBooks(NewBook); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}
	context.JSON(201, tools.OK)
	return
}

// DelVote 删除一个投票
func DelBooks(context *gin.Context) {
	idStr := context.Query("title")
	if err := model.DelBooks(idStr); err == nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code: 10006,
		})
		return
	}

	context.JSON(http.StatusOK, tools.OK)
	return
}

func BorrowBook(context *gin.Context) {
	name, _ := context.GetPostForm("title")
	jwt, _ := context.Cookie("jwt")
	JWT, err := model.CheckJwt(jwt)
	if err != nil {
		context.JSON(http.StatusOK, tools.BorrowErr2)
		return
	}
	userIDStr := JWT.Id
	userNameStr := JWT.Name
	if userIDStr < 0 {
		context.JSON(http.StatusOK, tools.BorrowErr2)
		return
	}
	if err := model.BorrowBook(userIDStr, userNameStr, name); err != nil {
		context.JSON(http.StatusOK, tools.BorrowErr)
		return
	}
	context.JSON(http.StatusOK, tools.OK1)
	return
}
func ReturnBook(context *gin.Context) {
	title, _ := context.GetPostForm("title")
	jwt, _ := context.Cookie("jwt")
	JWT, err := model.CheckJwt(jwt)
	if err != nil {
		context.JSON(http.StatusOK, tools.BorrowErr2)
	}
	userIDStr := JWT.Id
	if userIDStr < 0 {
		context.JSON(http.StatusOK, tools.BorrowErr2)
	}
	if err := model.ReturnBook(userIDStr, title); err != nil {
		context.JSON(http.StatusOK, tools.ReturnErr)
		return
	}

	context.JSON(http.StatusOK, tools.OK2)
	return
}
func GetImg(context *gin.Context) {
	var Img string
	title := context.Query("title")
	if img, err := model.GetImg(title); err != nil {
		context.JSON(http.StatusOK, tools.GetImgErr)
		return
	} else {
		Img = img
	}
	rootDir := "F:\\go.code\\src\\go_code\\lms\\app\\img\\" + Img
	fmt.Println(rootDir)
	imageData, err := os.ReadFile(rootDir)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read image",
		})
		return
	}
	// 设置响应头为图片类型
	context.Header("Content-Type", "image/jpeg")
	// 将图片数据作为响应主体发送给前端
	context.Data(http.StatusOK, "image/jpeg", imageData)

}
