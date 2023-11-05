package models

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func NewGormDB() {
	dsn := "root:123456789@tcp(localhost:3306)/admin_up?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// err = db.AutoMigrate(&UserBasic{}, &RoleBasic{}, &RoleMenu{}, &RoleFunction{}, &MenuBasic{}, &FunctionBasic{})
	// if err != nil {
	// 	panic(err)
	// }
	DB = db

}
func NewRedisDB() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RDB = rdb
}
