package app

import (
	"github.com/spf13/cobra"
	"go_code/lms/app/model"
	"go_code/lms/app/router"
	"log"
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
	rootCmd := &cobra.Command{
		Use:   "lms",
		Short: "Library Management System CLI",
		Run: func(cmd *cobra.Command, args []string) {
			router.New(cmd, args)
		},
	}
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
