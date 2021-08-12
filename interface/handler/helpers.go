package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func newCtx() context.Context {
	return context.Background()
}

func bind(c *gin.Context, request interface{}) (ok bool) {
	if err := c.BindJSON(request); err != nil {
		c.Status(http.StatusBadRequest)
		return false
	} else {
		return true
	}
}
