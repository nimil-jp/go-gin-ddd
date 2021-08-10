package persistence

import (
	"github.com/pkg/errors"
	"go-ddd/pkg/xerrors"
	"gorm.io/gorm"
)

func dbError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return xerrors.NotFound()
	default:
		return errors.WithStack(err)
	}
}
