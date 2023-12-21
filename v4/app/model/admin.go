package model

import (
	"fmt"
	"gorm.io/gorm"
)

func GetAdmin(name string) *Admin {
	var ret Admin
	if err := Conn.Table("admin").Where("name = ?", name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}
func GetAdminList(limit, offset int) []Admin {
	admin := make([]Admin, 0)
	if err := Conn.Table("admin").Offset(offset).Limit(limit).Find(&admin).Error; err != nil {
		fmt.Printf("图书管理员获取失败 err:%s", err.Error())
	}
	return admin
}
func GetInfosList(limit, offset int) []UserBooks {
	admin := make([]UserBooks, 0)
	if err := Conn.Table("user_books").Offset(offset).Limit(limit).Find(&admin).Error; err != nil {
		fmt.Printf("借阅记录获取失败 err:%s", err.Error())
	}
	return admin
}
func AddAdmin(admin Admin) error {
	//启用数据库事务添加新的表格
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&admin).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func UpdateAdmin(admin Admin) error {

	//启用数据库事务添加新的表格
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&admin).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func DelAdmin(admin string) error {

	//启用数据库事务添加新的表格
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(Admin{}, admin).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
