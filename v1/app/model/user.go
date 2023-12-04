package model

import (
	"fmt"
)

func GetUser(name string) *User {
	var ret User
	if err := Conn.Table("user").Where("name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

//func GetUserV1(name string) *User {
//	var ret User
//	err := Conn.Raw("select * from user where name = ? limit 1", name).Scan(&ret).Error
//	if err != nil {
//		tools.Logger.Errorf("[GetUserV1]err:%s", err.Error())
//	}
//	return &ret
//}

// CreateUser 参数是指针
func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}
	return nil
}
