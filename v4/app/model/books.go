package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func GetBook() []Books {
	books := make([]Books, 0)
	if err := Conn.Table("books").Find(&books).Error; err != nil {
		fmt.Printf("图书获取失败 err:%s", err.Error())
	}
	return books
}
func GetBookName(name string) Books {
	books := Books{}
	if err := Conn.Table("books").First(&books).Error; err != nil {
		fmt.Printf("图书获取失败 err:%s", err.Error())
	}
	return books
}
func GetBooks(limit, offset int) []Books {
	books := make([]Books, 0)
	if err := Conn.Table("books").Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		fmt.Printf("图书获取失败 err:%s", err.Error())
	}
	return books
}
func GetBooksUid(name string) int64 {
	user := User{}
	if err := Conn.Table("user").Where("name = ?", name).First(&user).Error; err == nil {
		return user.Uid
	}
	return 0
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
	//var num int
	var userBook UserBooks
	var Books Books
	//var existingStatus int
	fmt.Println(userId, userName, bookName)
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("books").Where("title = ?", bookName).Set("num", "FOR UPDATE").First(&Books).Error; err != nil {
			return errors.New("指定书名数量获取错误")
		}

		if Books.Num < 1 {
			return errors.New("库存不足")
		}
		if err := db.Table("user_books").Where("user_id = ? AND books_title = ?", userId, bookName).Last(&userBook).Error; err == nil {
			if userBook.Status == 1 {
				// 创建新的借阅记录
				mess := UserBooks{
					UserId:      userId,
					UserName:    userName,
					BooksTitle:  bookName,
					Status:      0, // 将状态设置为未还书
					CreatedTime: time.Now(),
					UpdatedTime: time.Now(),
				}
				if err := db.Create(&mess).Error; err != nil {
					fmt.Printf("err:%s", err.Error())
					return err
				}
				num := Books.Num - 1
				if err := db.Table("books").Where("id = ?", Books.Id).Update("num", num).Error; err != nil {
					return err
				}
				return nil
			} else if userBook.Status == 0 {
				return errors.New("你上次借的这本书还未归还")
			}
		} else {
			mess := UserBooks{
				UserId:      userId,
				UserName:    userName,
				BooksTitle:  bookName,
				Status:      0, // 将状态设置为未还书
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
			}
			if err := db.Create(&mess).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
			num := Books.Num - 1
			if err := db.Table("books").Where("id = ?", Books.Id).Update("num", num).Error; err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	return err

}

//	func ReturnBook(userId int64, bookName string) error {
//		var num int
//
//		err := Conn.Transaction(func(db *gorm.DB) error {
//			var userBook UserBooks
//			err := db.Table("user_books").Where("user_id = ? AND books_title = ?", userId, bookName).Order("created_time DESC").First(&userBook).Error
//			if errors.Is(err, gorm.ErrRecordNotFound) {
//				// 未找到借阅记录
//				return errors.New("未找到借阅记录")
//			} else if err != nil {
//				// 借阅记录查询异常
//				return errors.New("借阅记录查询异常")
//			}
//			if userBook.Status == 1 {
//				// 这本书已经归还
//				return errors.New("这本书已经归还")
//			} else if userBook.Status == 0 {
//				// 更新借阅记录状态为已归还
//				userBook.Status = 1
//				if err := db.Save(&userBook).Error; err != nil {
//					return err
//				}
//				// 还书成功后，num加一
//				num += 1
//				// 更新num
//				if err := db.Table("books").Where("title = ?", bookName).Update("num", num).Error; err != nil {
//					return err
//				}
//			}
//			return nil
//		})
//
//		return err
//	}
var (
	ErrRecordNotFound         = errors.New("未找到借阅记录")
	ErrBorrowRecordException  = errors.New("借阅记录查询异常")
	ErrBookAlreadyReturned    = errors.New("这本书已经归还")
	ErrOptimisticLockConflict = errors.New("乐观锁冲突")
)

func ReturnBook(userId int64, bookName string) error {
	err := Conn.Transaction(func(db *gorm.DB) error {
		var userBook UserBooks
		err := db.Table("user_books").Where("user_id = ? AND books_title = ?", userId, bookName).Order("created_time DESC").First(&userBook).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未找到借阅记录
			return ErrRecordNotFound
		} else if err != nil {
			// 借阅记录查询异常
			return ErrBorrowRecordException
		}
		if userBook.Status == 1 {
			// 这本书已经归还
			return ErrBookAlreadyReturned
		} else if userBook.Status == 0 {
			// 更新借阅记录状态为已归还
			userBook.Status = 1
			userBook.Version++ // 增加版本号

			if err := db.Save(&userBook).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 乐观锁冲突
					return ErrRecordNotFound
				}
				return err
			}

			// 还书成功后，num加一
			var book Books
			if err := db.Table("books").Where("title = ?", bookName).First(&book).Error; err != nil {
				return err
			}
			book.Num += 1 // 增加数量

			if err := db.Save(&book).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
func GetImg(title string) (string, error) {
	var book Books
	fmt.Println(title)
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("books").Where("title = ?", title).Select("img_url").First(&book).Error; err != nil {
			fmt.Println()
			return err
		}
		if book.ImgUrl == "" {
			return errors.New("未找到匹配的书籍")
		}
		return nil
	})
	fmt.Println(book.ImgUrl)
	return book.ImgUrl, err
}
