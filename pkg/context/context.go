package context

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-gin-ddd/pkg/xerrors"
)

type Context interface {
	RequestID() string
	UserID() uint

	Validate(request interface{}) (ok bool)
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
	userID uint
}

func New(requestID string, userID uint) Context {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return &ctx{
		id:     requestID,
		verr:   xerrors.NewValidation(),
		userID: userID,
	}
}

func (c ctx) RequestID() string {
	return c.id
}

func (c ctx) UserID() uint {
	return c.userID
}
