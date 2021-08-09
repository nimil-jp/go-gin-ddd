package xerrors

import (
	"fmt"
	"net/http"
)

type Expected struct {
	statusCode int
	msg        string
}

func NewExpected(statusCode int, message string) Expected {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	return Expected{statusCode: statusCode, msg: message}
}

func (e Expected) StatusCode() int {
	return e.statusCode
}

func (e Expected) Message() string {
	return e.msg
}

func (e Expected) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.statusCode, e.msg)
}
