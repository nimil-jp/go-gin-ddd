package main

import (
	"go-ddd/domain/entity"
	"go-ddd/util"
)

func init() {
	err := util.DB.AutoMigrate(
		&entity.User{},
	)
	if err != nil {
		panic(err)
	}
}
