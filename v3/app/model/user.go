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
func GetUser2(name string) *User {
	var ret User
	if err := Conn.Table("user").Where("email = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}
func GetUser3(telephone string) *User {
	var ret User
	if err := Conn.Table("user").Where("telephone = ?", telephone).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

// CreateUser 参数是指针
func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}
	return nil
}
func GetUserImg() {

}
