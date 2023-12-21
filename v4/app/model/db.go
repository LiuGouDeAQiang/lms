package model

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// 数据库操作都放在这里
var MB *mongo.Client
var Conn *gorm.DB
var Rdb *redis.Client

func NewMongoDB() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://192.168.10.29:27017")
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	MB = client

	fmt.Println("Connected to MongoDB!")
}
func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "123456", "127.0.0.1:3306", "lms")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	Conn = conn
}

func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.10.29:6379",
		Password: "",
		DB:       0,
	})
	Rdb = rdb
	//初始化session
	Store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)
	return
}

//用户 --> 客户端应用程序: 输入手机号和验证码
//客户端应用程序 --> 服务器: 发送手机号和验证码
//服务器 --> 客户端应用程序: 验证手机号和验证码
//客户端应用程序 --> 用户: 登录成功
// 邮箱验证码部分

func NewEmail(email string) (string, error) {
	smtpHost := "smtp.qq.com"
	smtpPort := 587
	sender := "3494383150@qq.com"  // 请将此处替换为您的有效Gmail邮箱地址
	password := "pjgdlqpjtnsucihb" // 请将此处替换为您的Gmail邮箱密码

	userTelephone := email
	rand.NewSource(time.Now().UnixNano())
	verificationCode := strconv.Itoa(rand.Intn(900000) + 100000)

	// 连接到SMTP服务器
	auth := smtp.PlainAuth("", sender, password, smtpHost)
	smtpAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	// 使用TLS加密连接
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 将此设置为false以启用证书验证
		ServerName:         smtpHost,
	}
	// 建立与SMTP服务器的连接
	client, err := smtp.Dial(smtpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 开始TLS握手
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Fatal(err)
	}

	// 使用认证信息进行登录
	if err = client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// 设置发件人
	if err = client.Mail(sender); err != nil {
		log.Fatal(err)
	}

	// 设置收件人
	if err = client.Rcpt(userTelephone); err != nil {
		log.Fatal(err)
	}

	// 获取写入流
	w, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	// 构建邮件内容
	subject := "Verification Code"
	body := "Your verification code is: " + verificationCode
	msg := []string{
		"To: " + userTelephone,
		"From: " + sender,
		"Subject: " + subject,
		"",
		body,
	}

	// 写入邮件内容
	_, err = w.Write([]byte(fmt.Sprintf("%s\r\n", strings.Join(msg, "\r\n"))))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Verification code:", verificationCode)
	return verificationCode, err
}

const (
	RegionID        = "cn-hangzhou"
	AccessKeyID     = "LTAI5tL6sVPQmoBMX9xR5KQz"
	AccessKeySecret = "MpG3MQAKNWJRddSauIEBB4F7u0HHWw"
	SignName        = "阿里云短信测试"
	TemplateCode    = "SMS_154950909"
	TemplateParam   = `{"code":"123456"}`
	PhoneNumber     = "YourPhoneNumber"
)

// 短信验证部分

func GenerateVerificationCode() string {
	rand.NewSource(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
func NewTelephone(phoneNumber string) (string, error) {
	client, err := dysmsapi.NewClientWithAccessKey(RegionID, AccessKeyID, AccessKeySecret)
	if err != nil {
		return "", err
	}
	verificationCode := GenerateVerificationCode()
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumber
	request.SignName = "阿里云短信测试"
	request.TemplateCode = "SMS_154950909"
	request.TemplateParam = `{"code":"` + verificationCode + `"}`

	response, err := client.SendSms(request)
	if err != nil {
		return "", err
	}
	if response.Code != "OK" {
		return "", errors.New(response.Message)
	}

	return verificationCode, nil
}

func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
	_ = Rdb.Close()
}

//用户 --> 客户端应用程序: 打开登录页面
//客户端应用程序 --> 微信服务器: 生成登录二维码
//微信服务器 --> 客户端应用程序: 返回登录二维码
//客户端应用程序 --> 用户: 显示登录二维码
//用户 --> 微信客户端: 扫描登录二维码
//微信客户端 --> 微信服务器: 发送扫描请求
//微信服务器 --> 微信客户端: 返回扫描状态（未登录/已登录）
//微信客户端 --> 客户端应用程序: 返回扫描状态
//客户端应用程序 --> 微信服务器: 发送登录确认请求
//微信服务器 --> 客户端应用程序: 返回登录成功
//客户端应用程序 --> 用户: 登录成功
