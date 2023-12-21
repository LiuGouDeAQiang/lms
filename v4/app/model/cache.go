package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"

	"time"
)

func GetBooksFromCache(key string) ([]Books, error) {
	// 尝试从缓存中获取数据
	val, err := Rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// 缓存中没有数据
		return nil, fmt.Errorf("cache miss")
	} else if err != nil {
		// 发生其他错误
		return nil, err
	}

	// 解码缓存数据
	var books []Books
	err = json.Unmarshal([]byte(val), &books)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// SetBooksToCache 将图书列表存储到string  Redis缓存中
func SetBooksToCache(key string, books []Books) error {
	// 编码图书列表为JSON
	data, err := json.Marshal(books)
	if err != nil {
		return err
	}
	// 将数据存储到缓存中，设置过期时间为1小时
	err = Rdb.Set(context.Background(), key, data, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetBooksToZSetCache 将图书列表存储到 Redis zset 类型的数据结构中
func SetBooksToZSetCache(key string, books []Books) error {
	pipe := Rdb.Pipeline()

	// 存储到 Redis zset 类型的数据结构中，使用当前时间作为分数
	for _, book := range books {
		data, err := json.Marshal(book)
		if err != nil {
			return err
		}
		member := string(data) // 将字节切片转换为字符串类型

		pipe.ZAdd(context.Background(), key, redis.Z{Score: float64(time.Now().UnixNano()), Member: member})
	}

	_, err := pipe.Exec(context.Background())
	return err
}

// SetBooksToHashCache 将图书列表存储到 Redis hash 类型的数据结构中
func SetBooksToHashCache(key string, books []Books) error {
	// 存储到 Redis hash 类型的数据结构中，将图书列表的每本书存储为 hash 的一个 field-value 对
	pipe := Rdb.Pipeline()
	for _, book := range books {
		data, err := json.Marshal(book)
		if err != nil {
			return err
		}
		pipe.HSet(context.Background(), key, book.Title, data)
	}
	_, err := pipe.Exec(context.Background())
	return err
}

func GetAdminsFromCache(key string) ([]Admin, error) {
	// 尝试从缓存中获取数据
	val, err := Rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// 缓存中没有数据
		return nil, fmt.Errorf("查询为空")
	} else if err != nil {
		// 发生其他错误
		return nil, err
	}
	// 解码缓存数据
	var admin []Admin
	err = json.Unmarshal([]byte(val), &admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// SetAdminToCache SetBooksToCache 将图书列表存储到string  Redis缓存中
func SetAdminToCache(key string, admin []Admin) error {
	// 编码图书列表为JSON
	data, err := json.Marshal(admin)
	if err != nil {
		return err
	}
	// 将数据存储到缓存中，设置过期时间为1小时
	err = Rdb.Set(context.Background(), key, data, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetInfosFromCache(key string) ([]UserBooks, error) {
	// 尝试从缓存中获取数据
	val, err := Rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		// 缓存中没有数据
		return nil, fmt.Errorf("查询为空")
	} else if err != nil {
		// 发生其他错误
		return nil, err
	}
	// 解码缓存数据
	var admin []UserBooks
	err = json.Unmarshal([]byte(val), &admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
func SetInfosToCache(key string, admin []UserBooks) error {
	// 编码图书列表为JSON
	data, err := json.Marshal(admin)
	if err != nil {
		return err
	}
	// 将数据存储到缓存中，设置过期时间为1小时
	err = Rdb.Set(context.Background(), key, data, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}
