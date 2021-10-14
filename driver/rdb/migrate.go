package rdb

import "go-gin-ddd/domain/entity"

func migrate() error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
