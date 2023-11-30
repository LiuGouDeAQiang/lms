package model

import (
	"context"
	"fmt"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库操作都放在这里

var Conn *gorm.DB
var Rdb *redis.Client

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

func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
	_ = Rdb.Close()
}
