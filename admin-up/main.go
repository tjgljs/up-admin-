package main

import (
	"admin-up/models"
	"admin-up/router"
)

func main() {
	// 初始化 gorm.DB
	models.NewGormDB()
	models.NewRedisDB()

	r := router.App()

	r.Run()
}
