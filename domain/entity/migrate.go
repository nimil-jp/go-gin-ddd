package entity

import "go-ddd/util"

func init() {
	err := util.DB.AutoMigrate(
		&User{},
	)
	if err != nil {
		panic(err)
	}
}
