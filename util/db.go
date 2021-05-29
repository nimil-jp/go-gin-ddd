package util

import (
	"fmt"

	"go-ddd/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	var con string

	if config.Env.DB.Socket != "" {
		con = fmt.Sprintf("unix(%s)", config.Env.DB.Socket)
	} else {
		con = fmt.Sprintf("tcp(%s:%d)", config.Env.DB.Host, config.Env.DB.Port)
	}

	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DB.User,
		config.Env.DB.Password,
		con,
		config.Env.DB.Name,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

}
