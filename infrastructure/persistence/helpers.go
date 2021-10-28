package persistence

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/nimil-jp/gin-utils/xerrors"
)

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return xerrors.NotFound()
	default:
		return errors.WithStack(err)
	}
}
