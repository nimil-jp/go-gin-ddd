package context

import (
	"gorm.io/gorm"

	"go-gin-ddd/driver/rdb"
)

func (c *ctx) DB() *gorm.DB {
	if c.db != nil {
		return c.db
	}
	return rdb.Get()
}

func (c *ctx) Transaction(fn func(ctx Context) error) error {
	return rdb.Get().Transaction(
		func(tx *gorm.DB) error {
			c.db = tx

			return fn(c)
		},
	)
}
