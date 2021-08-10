package xerrors

import (
	"fmt"
	"net/http"
)

type Expected struct {
	statusCode int
	msg        string
}

func NewExpected(statusCode int, message string) *Expected {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return &Expected{statusCode: statusCode, msg: message}
}

func (e Expected) StatusCode() int {
	return e.statusCode
}

func (e Expected) Message() string {
	return e.msg
}

func (e *Expected) ChangeStatus(before int, after int) bool {
	if e.statusCode == before {
		e.statusCode = after
		return true
	}
	return false
}

func (e Expected) StatusOk() bool {
	return e.statusCode < 300
}

func (e Expected) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.statusCode, e.msg)
}

// expected errors

func NotFound() *Expected {
	return NewExpected(http.StatusNotFound, "not found")
}
