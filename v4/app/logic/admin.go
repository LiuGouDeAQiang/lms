package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/model"
	"go_code/lms/app/tools"
	"net/http"
	"strconv"
	"time"
)

type Admin struct {
	Name         string `json:"name" form:"name"`
	RoleId       string `json:"role_id" form:"role_id"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetAdminLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "adminLogin.html", nil)
}
func GetDOAdminLogin(context *gin.Context) {
	var admin Admin
	if err := context.ShouldBind(&admin); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(), //这里有风险
		})
	}
	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: admin.CaptchaId,
		Data:      admin.CaptchaValue,
	}) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}
	ret := model.GetAdmin(admin.Name)
	fmt.Println(ret)
	if ret.Id < 1 || ret.Password != admin.Password {
		context.JSON(http.StatusOK, tools.UserErr)
		return
	}
	_ = model.SetSession(context, admin.Name, ret.Id, ret.RoleId)
	nuy, err := model.GetJwt(ret.Id, admin.Name)
	context.SetCookie("jwt", nuy, 3600, "/", "", true, false)
	if err != nil {
		panic(err)
	}
	context.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
	})
	return
}
func GetAdminList0(c *gin.Context) {
	c.HTML(http.StatusOK, "adminList.html", nil)
}
func GetAdminList(c *gin.Context) {
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
	cacheKey := fmt.Sprintf("admins:%d:%d", limit, offset)
	//尝试从Redis缓存中获取管理员信息
	admins, err := model.GetAdminsFromCache(cacheKey)
	if err == nil {
		// 在缓存中找到数据，直接返回缓存结果
		c.JSON(http.StatusOK, gin.H{"admins": admins})
		return
	}
	// 调用函数获取图书列表（分页）
	admins = model.GetAdminList(limit, offset)
	// 将图书列表存储到Redis缓存中
	fmt.Println(admins)
	err = model.SetAdminToCache(cacheKey, admins)
	if err != nil {
		// 处理存储到缓存中的错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store books in cache"})
		return
	}
	// 返回图书列表
	c.JSON(http.StatusOK, gin.H{"admins": admins})
}

func AdminAdd(context *gin.Context) {
	name, _ := context.GetPostForm("name")
	password, _ := context.GetPostForm("password")
	uid, _ := context.GetPostForm("uid")
	UID, err := strconv.ParseInt(uid, 10, 32)
	if err != nil {
		panic(err)
	}
	//构建结构体
	NewAdmin := model.Admin{
		Name:        name,
		Password:    password,
		Uid:         UID,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}
	if NewAdmin.Name == " " {
		context.JSON(http.StatusBadRequest, tools.ECode{
			Message: "创建失败，参数错误",
			Data:    nil,
		})
	}
	if err := model.AddAdmin(NewAdmin); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusCreated, tools.OK)
	return
}
func AdminUpdate(context *gin.Context) {
	id, _ := context.GetPostForm("id")
	name, _ := context.GetPostForm("name")
	password, _ := context.GetPostForm("password")
	uid, _ := context.GetPostForm("uid")
	UID, err := strconv.ParseInt(uid, 10, 32)
	Id, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		panic(err)
	}
	//构建结构体
	NewAdmin := model.Admin{
		Id:          Id,
		Name:        name,
		Password:    password,
		Uid:         UID,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}
	if NewAdmin.Name == " " {
		context.JSON(http.StatusBadRequest, tools.ECode{
			Message: "修改失败，参数错误",
			Data:    nil,
		})
	}
	if err := model.UpdateAdmin(NewAdmin); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, tools.OK)
	return
}
func AdminDel(context *gin.Context) {
	uid, _ := context.GetPostForm("uid")
	if err := model.DelAdmin(uid); err == nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code: 10006,
		})
		return
	}
	context.JSON(http.StatusOK, tools.OK)
	return
}
func GetInfos0(c *gin.Context) {
	c.HTML(http.StatusOK, "borrowingInfos.html", nil)
}
func GetInfos(c *gin.Context) {
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
	cacheKey := fmt.Sprintf("user_books:%d:%d", limit, offset)
	//尝试从Redis缓存中获取管理员信息na
	infos, err := model.GetInfosFromCache(cacheKey)
	if err == nil {
		// 在缓存中找到数据，直接返回缓存结果
		c.JSON(http.StatusOK, gin.H{"infos": infos})
		return
	}
	// 调用函数获取图书列表（分页）
	infos = model.GetInfosList(limit, offset)
	// 将图书列表存储到Redis缓存中
	err = model.SetInfosToCache(cacheKey, infos)
	if err != nil {
		// 处理存储到缓存中的错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store books in cache"})
		return
	}
	// 返回图书列表
	c.JSON(http.StatusOK, gin.H{"infos": infos})
}

func GetAdminView(context *gin.Context) {
	context.HTML(http.StatusOK, "adminView.html", nil)
}
