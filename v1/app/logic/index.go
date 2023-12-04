package logic

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/model"
	"go_code/lms/app/tools"
	"net/http"
	"strconv"
	"time"
)

func Index(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")   // 设置默认值为5
	offsetStr := c.DefaultQuery("offset", "5") // 设置默认值为5

	// 将分页参数转换为整数
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		// 处理分页参数错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		// 处理分页参数错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	// 调用函数获取图书列表（分页）
	books := model.GetBooks(limit, offset)
	if err != nil {
		// 处理获取图书列表错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books"})
		return
	}

	// 返回图书列表
	c.HTML(http.StatusOK, "index.html", gin.H{"books": books})
}
func GetBooks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")   // 设置默认值为5
	offsetStr := c.DefaultQuery("offset", "5") // 设置默认值为5

	// 将分页参数转换为整数
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		// 处理分页参数错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		// 处理分页参数错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	// 调用函数获取图书列表（分页）
	books := model.GetBooks(limit, offset)
	if err != nil {
		// 处理获取图书列表错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get books"})
		return
	}

	// 返回图书列表
	c.JSON(http.StatusOK, gin.H{"books": books})
}

// CheckXYZ 限流
func CheckXYZ(context *gin.Context) bool {
	ip := context.ClientIP()
	ua := context.GetHeader("user-agent")
	hash := md5.New()           //创建一个MD5哈希实例
	hash.Write([]byte(ip + ua)) //将IP地址和user-agent信息拼接后写入哈希实例。
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes) //将哈希值转换为字符串

	flag, _ := model.Rdb.Get(context, "ban-"+hashString).Bool()
	if flag {
		return false
	}
	i, _ := model.Rdb.Get(context, "xyz-"+hashString).Int() // 从Redis中获取"xyz-"+hashString键对应的值，并将其转换为整数类型
	if i > 10 {
		model.Rdb.SetEx(context, "ban-"+hashString, true, 3*time.Second)
		return false
	}
	// 如果获取的值大于5，则将"ban-"+hashString键设置为true（加入黑名单），并设置过期时间为3秒
	model.Rdb.Incr(context, "xyz-"+hashString)                  //Incr将存储值递增一，Expire用于设置过期时间
	model.Rdb.Expire(context, "xyz-"+hashString, 5*time.Second) //每次访问时次数加一，并设置过期时间5秒
	return true
}
func GetCaptcha(context *gin.Context) {
	if !CheckXYZ(context) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10016,
			Message: "限流了",
		})
		return
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, tools.ECode{
		Data: captcha,
	})
}
