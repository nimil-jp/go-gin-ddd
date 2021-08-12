package tx

import (
	"context"

	"go-ddd/util"
	"gorm.io/gorm"
)

var key = struct{}{}

type txFunc func(ctx context.Context) error

func Do(ctx context.Context, f txFunc) error {
	return util.DB.Transaction(
		func(tx *gorm.DB) error {
			ctx := context.WithValue(ctx, &key, tx)

			return f(ctx)
		},
	)
}

func Get(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(&key).(*gorm.DB)
	if ok {
		return tx
	}
	return util.DB
}
