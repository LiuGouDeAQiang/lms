package model

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

type UserToken struct {
	Id   int64
	Name string
	jwt.RegisteredClaims
}

// 签名密钥

func GetJwt(id int64, name string) (string, error) {
	if id < 0 || name == "" {
		return "", errors.New("参数错误")
	}
	//读取签名密钥配置
	signKey := viper.GetString("jwt.signKey")
	token := &UserToken{
		Id:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "可求帅图书馆",                                             // 签发者
			Subject:   "名流张三",                                               // 签发对象
			Audience:  jwt.ClaimStrings{"Android", "IOS", "H5"},             //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),        //过期时间 1小时
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * 10)), //最早使用时间 10秒之后
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       //签发时间 当前时间
			ID:        "Test-1",                                             // jwt ID,类似于盐值 最好是每次都随机
		},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(signKey))
	return tokenStr, err
}

func CheckJwt(tokenStr string) (*UserToken, error) {
	signKey := viper.GetString("jwt.signKey")
	token, err := jwt.ParseWithClaims(tokenStr, &UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil //返回签名密钥
	})
	if err != nil || !token.Valid {
		return nil, errors.New("校验失败，TOKEN不合格")
	}

	claims, ok := token.Claims.(*UserToken)
	if !ok {
		return nil, errors.New("TOKEN转义失败！")
	}

	return claims, nil
}
