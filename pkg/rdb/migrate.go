package rdb

import "go-ddd/domain/entity"

func migrate() error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
