package app

import (
	"go_code/lms/app/model"
	"go_code/lms/app/router"
)

// Start 启动器方法
func Start() {
	model.NewMysql()
	model.NewRdb()
	defer func() {
		model.Close()
	}()
	//定时器
	//schedule.Start()
	//tools.NewLogger()
	router.New()
}
