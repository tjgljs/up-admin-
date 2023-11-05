package test

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDB(t *testing.T) {
	dsn := "root:123456789@tcp(localhost:3306)/admin_up?charset=utf8&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

}
