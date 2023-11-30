package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetBooks() []Books {
	books := make([]Books, 0)
	if err := Conn.Table("books").Find(&books).Error; err != nil {
		fmt.Printf("图书获取失败 err:%s", err.Error())
	}

	return books
}

func AddBooks(book Books) error {

	//启用数据库事务添加新的表格
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&book).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func UpdateBooks(books Books) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&books).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func DelBooks(title string) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Books{}, title).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func BorrowBook(userId int64, userName string, bookName string) error {
	var num int
	var existingStatus int
	fmt.Println(userId, userName, bookName)
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("books").Where("title = ?", bookName).Select("num").Scan(&num).Error; err != nil {
			return err
		}
		if err := db.Table("user_books").Where("user_id = ? AND books_title = ?", userId, bookName).Select("status").Scan(&existingStatus).Error; err != nil {
			return err
		}

		if existingStatus == 0 {
			return errors.New("你已经借过这本书了")
		}
		if num >= 1 {
			num-- // 减一
			if err := db.Table("books").Where("title = ?", bookName).Update("num", num).Error; err != nil {
				return err
			}
			mess := UserBooks{
				UserId:      userId,
				UserName:    userName,
				BooksTitle:  bookName,
				Status:      0,
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
			}
			if err := db.Create(&mess).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
			return nil
		}
		return errors.New("图书库存不足")
	})
	return err
}
func ReturnBook(userId int64, bookName string) error {
	var num int
	var status int

	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("books").Where("title = ?", bookName).Select("num").Scan(&num).Error; err != nil {
			return err
		}

		if err := db.Table("user_books").Where("user_id = ? AND books_title = ?", userId, bookName).Select("status").Scan(&status).Error; err != nil {
			return err
		}

		if status == 0 {
			num++ // 增加一本书
			if err := db.Table("books").Where("title = ?", bookName).Update("num", num).Error; err != nil {
				return err
			}

			if err := db.Table("user_books").Where("user_id = ? AND books_title = ?", userId, bookName).Update("status", 1).Error; err != nil {
				return err
			}

			return nil
		} else if status == 1 {
			return errors.New("你已经还过这本书了")
		}

		return errors.New("无效的借阅记录")
	})

	return err
}
