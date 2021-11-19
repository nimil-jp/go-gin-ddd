package persistence

import (
	"gorm.io/gorm"

	"github.com/nimil-jp/gin-utils/errors"
)

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.NotFound()
	default:
		return errors.NewUnexpected(err)
	}
}
