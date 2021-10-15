package context

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-gin-ddd/pkg/xerrors"
)

type Context interface {
	RequestId() string
	UserId() uint

	FieldError(fieldName string, message string)
	IsInValid() bool
	ValidationError() error

	DB() *gorm.DB
	Transaction(fn func(ctx Context) error) error
}

type ctx struct {
	id     string
	verr   *xerrors.Validation
	db     *gorm.DB
	userId uint
}

func New(requestId string, userId uint) Context {
	if requestId == "" {
		requestId = uuid.New().String()
	}
	return &ctx{
		id:     requestId,
		verr:   xerrors.NewValidation(),
		userId: userId,
	}
}

func (c ctx) RequestId() string {
	return c.id
}

func (c ctx) UserId() uint {
	return c.userId
}
