package logic

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetUserImg(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 创建存储目录
	err = os.MkdirAll("F:\\go.code\\src\\go_code\\lms\\app\\user_img", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create upload directory")
	}
	// 生成保存文件的路径
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join("F:\\go.code\\src\\go_code\\lms\\app\\user_img\\", filename)

	// 将文件保存到本地
	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file":    savePath,
	})

}
