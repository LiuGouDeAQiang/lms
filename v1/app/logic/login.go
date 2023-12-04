package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"go_code/lms/app/model"

	"go_code/lms/app/tools"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

//用户                          系统
//|                             |
//|---------发送登录请求--------->|
//|                             |
//|<-------返回登录页面----------|
//|                             |
//|---------提交登录凭证--------->|
//|                             |
//|<-------验证登录凭证----------|
//|                             |
//|--------返回登录成功---------->|
//|                             |
//

type User struct {
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetLogin(context *gin.Context) {
	context.HTML(http.StatusOK, "login.html", nil)
}
func GetCreate(context *gin.Context) {
	context.HTML(http.StatusOK, "create.html", nil)
}
func DoLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: err.Error(), //这里有风险
		})
	}
	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败！", //这里有风险
		})
		return
	}
	//数据库查询
	ret := model.GetUser(user.Name)
	if ret.Id < 1 /*|| ret.Password != user.Password*/ {
		context.JSON(http.StatusOK, tools.UserErr)
		return
	}
	_ = model.SetSession(context, user.Name, ret.Id)
	fmt.Printf(user.Name, ret.Id)
	nuy, err := model.GetJwt(ret.Id, user.Name)
	context.SetCookie("jwt", nuy, 3600, "/", "", true, false)
	if err != nil {
		panic(err)
	}
	context.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
	})
	return
}

func Logout(context *gin.Context) {

	_ = model.FlushSession(context)
	context.Redirect(http.StatusFound, "/login")
}

type CUser struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
	Telephone string `json:"phone"`
}

func CreateUser(context *gin.Context) {
	var user CUser
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(), //这里有风险
		})
		return
	}
	fmt.Println(user.Name, user.Password, user.Password2, user.Telephone)
	//encryptV1(user.Password)
	//对数据进行校验
	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "账号或者密码不能为空", //这里有风险
		})
		return
	}

	//校验密码
	if user.Password != user.Password2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同！", //这里有风险
		})
		return
	}
	//校验用户是否存在，这种写法非常不安全。有严重地并发风险
	if oldUser := model.GetUser(user.Name); oldUser.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在", //这里有风险
		})
		return
	}
	//判断位数
	lenName := len(user.Name)
	lenPwd := len(user.Password)
	if lenName < 8 || lenName > 16 || lenPwd < 8 || lenPwd > 16 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "用户名或者密码要大于等于8，小于等于16！", //这里有风险
		})
		return
	}

	//密码不能是纯数字
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.Password) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字", //这里有风险
		})
		return
	}
	//开始添加用户
	newUser := model.User{
		Name:        encrypt(user.Name),
		Password:    encryptV1(user.Password),
		Telephone:   encrypt(user.Telephone),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Uid:         tools.GetUid(),
	}
	if err := model.CreateUser(&newUser); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "用户创建失败", //这里有风险
		})
		return
	}

	//返回添加成功
	context.JSON(http.StatusOK, tools.OK)
	return
}
func encrypt(pwd string) string {
	//创建md5hash.Has
	hash := md5.New()
	//将pwd写入MD5哈希对象
	hash.Write([]byte(pwd))
	//计算哈希值并返回字节数组
	hashBytes := hash.Sum(nil)
	//将字节数组转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码：%s\n", hashString)

	return hashString
}

func encryptV1(pwd string) string {
	//将原始密码与固定字符串连接起来，增加密码的复杂度
	newPwd := pwd + "可求帅图书馆" //不能随便起，且不能暴露
	hash := md5.New()
	hash.Write([]byte(newPwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码：%s\n", hashString)

	return hashString
}

func encryptV2(pwd string) string {
	//基于Blowfish 实现加密。简单快速，但有安全风险
	//golang.org/x/crypto/ 中有大量的加密算法
	newPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("密码加密失败：", err)
		return ""
	}
	newPwdStr := string(newPwd)
	fmt.Printf("加密后的密码：%s\n", newPwdStr)
	return newPwdStr
}
