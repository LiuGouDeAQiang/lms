package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rbcervilla/redisstore/v9"
	"net/http"
)

var Store *redisstore.RedisStore

// 代表会话储存，并进行一个新实例初始化，并采用字节切片作为参数
var SessionName = "session-name"

func GetSession(c *gin.Context) map[interface{}]interface{} {
	//从store存储中检索sessionName
	session, _ := Store.Get(c.Request, SessionName)
	//打印值
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}
func SetSession(c *gin.Context, name string, id int64) error {
	session, _ := Store.Get(c.Request, SessionName)
	session.Options.MaxAge = 3600
	session.Values["name"] = name
	session.Values["id"] = id
	//c.Request是当前请求的*http.Request对象，
	//c.Writer是当前请求的http.ResponseWriter对象。
	//这两个对象提供了对当前HTTP请求和响应的访问和控制。
	return session.Save(c.Request, c.Writer)
}
func FlushSession(c *gin.Context) error {
	session, _ := Store.Get(c.Request, SessionName)

	// Check if the user is a guest (not logged in)
	if session.Values["name"] == nil || session.Values["name"].(string) == "" {
		// Redirect to the login page
		c.Redirect(http.StatusFound, "/login")
		return nil
	}
	// Clear the session values
	session.Values["name"] = ""
	session.Values["id"] = int64(0)

	// Save the session to update the changes
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		return err
	}

	// Clear the cookies
	c.SetCookie("name", "", -1, "/", "", true, false)
	c.SetCookie("Id", "", -1, "/", "", true, false)

	return nil
}
