package rdb

import (
	"context"

	"gorm.io/gorm"
)

var key = struct{}{}

func Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return db.Transaction(
		func(tx *gorm.DB) error {
			ctx := context.WithValue(ctx, &key, tx)

			return fn(ctx)
		},
	)
}

func Get(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(&key).(*gorm.DB)
	if ok {
		return tx
	}
	return db
}
