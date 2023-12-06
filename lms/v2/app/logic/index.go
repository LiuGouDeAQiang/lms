package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/model"
	"go_code/lms/app/tools"
	"net/http"
	"strconv"
	"time"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GetBooks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")   // 设置默认值为5
	offsetStr := c.DefaultQuery("offset", "0") // 设置默认值为5
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
	cacheKey := fmt.Sprintf("books:%d:%d", limit, offset)
	// 尝试从Redis缓存中获取图书列表
	books, err := model.GetBooksFromCache(cacheKey)
	if err == nil {
		// 在缓存中找到数据，直接返回缓存结果
		c.JSON(http.StatusOK, gin.H{"books": books})
		return
	}
	// 调用函数获取图书列表（分页）
	books = model.GetBooks(limit, offset)
	// 将图书列表存储到Redis缓存中
	err = model.SetBooksToCache(cacheKey, books)
	if err != nil {
		// 处理存储到缓存中的错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store books in cache"})
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
	fmt.Println(i)
	return true
}
func XYZRateLimitMiddleware(c *gin.Context) bool {
	var name string
	values := model.GetSession(c)
	if v, ok := values["name"]; ok {
		name = v.(string)
		uid := strconv.FormatInt(model.GetBooksUid(name), 10)
		uidBytes := []byte(uid)
		fmt.Println(name, uid)
		hash := md5.New()    //创建一个MD5哈希实例
		hash.Write(uidBytes) //将IP地址和user-agent信息拼接后写入哈希实例。
		hashBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashBytes) //将哈希值转换为字符串
		flag, _ := model.Rdb.Get(c, "ban-"+hashString).Bool()
		if flag {
			return false
		}
		i, _ := model.Rdb.Get(c, "xyz-"+hashString).Int() // 从Redis中获取"xyz-"+hashString键对应的值，并将其转换为整数类型
		if i > 5 {
			model.Rdb.SetEx(c, "ban-"+hashString, true, 30*time.Second)
			return false
		}
		// 如果获取的值大于5，则将"ban-"+hashString键设置为true（加入黑名单），并设置过期时间为3秒
		model.Rdb.Incr(c, "xyz-"+hashString)                   //Incr将存储值递增一，Expire用于设置过期时间
		model.Rdb.Expire(c, "xyz-"+hashString, 50*time.Second) //每次访问时次数加一，并设置过期时间5秒
		return true
	} else {
		ip := c.ClientIP()
		ua := c.GetHeader("user-agent")
		fmt.Println(ip, ua)
		hash := md5.New()           //创建一个MD5哈希实例
		hash.Write([]byte(ip + ua)) //将IP地址和user-agent信息拼接后写入哈希实例。
		hashBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashBytes) //将哈希值转换为字符串

		flag, _ := model.Rdb.Get(c, "ban-"+hashString).Bool()
		if flag {
			return false
		}
		i, _ := model.Rdb.Get(c, "xyz-"+hashString).Int() // 从Redis中获取"xyz-"+hashString键对应的值，并将其转换为整数类型
		if i > 5 {
			model.Rdb.SetEx(c, "ban-"+hashString, true, 10*time.Second)
			return true
		}
		// 如果获取的值大于5，则将"ban-"+hashString键设置为true（加入黑名单），并设置过期时间为3秒
		model.Rdb.Incr(c, "xyz-"+hashString)                  //Incr将存储值递增一，Expire用于设置过期时间
		model.Rdb.Expire(c, "xyz-"+hashString, 5*time.Second) //每次访问时次数加一，并设置过期时间5秒
		return true
	}

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
func Xyz(context *gin.Context) {
	if !XYZRateLimitMiddleware(context) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    1644,
			Message: "要死啊，点那麽多",
		})
		return
	}
}
