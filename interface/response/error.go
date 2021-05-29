package response

import (
	"fmt"
	"log"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Error struct {
	statusCode int
	msg        string
}

func NewError(statusCode int, err error) Error {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return Error{statusCode: statusCode, msg: err.Error()}
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.statusCode, e.msg)
}

func ErrorJSON(c *gin.Context, err error) {
	var handleError Error
	if errors.As(err, &handleError) {
		c.JSON(handleError.statusCode, handleError.msg)
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
