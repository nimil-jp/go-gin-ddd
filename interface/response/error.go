package response

import (
	"fmt"
	"log"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"go-ddd/util/validate"
	"gorm.io/gorm"
)

type ExpectedError struct {
	statusCode int
	msg        string
}

func NewExpectedError(statusCode int, message string) ExpectedError {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return ExpectedError{statusCode: statusCode, msg: message}
}

func (e ExpectedError) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.statusCode, e.msg)
}

func ErrorJSON(c *gin.Context, err error) {
	var (
		eerr ExpectedError
		verr validate.ValidationError
	)

	if errors.As(err, &eerr) {
		c.JSON(eerr.statusCode, eerr.msg)
	} else if errors.As(err, &verr) {
		c.JSON(http.StatusBadRequest, verr)
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errors.New("record not found"))
		} else {
			if gin.Mode() == gin.DebugMode {
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
	}
	log.Printf("%+v\n", err)
}
