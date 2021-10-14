package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"go-gin-ddd/pkg/context"
	"go-gin-ddd/pkg/xerrors"
)

type handlerFunc func(ctx context.Context, c *gin.Context) error

func get(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.GET(relativePath, hf(handlerFunc))
}

func post(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.POST(relativePath, hf(handlerFunc))
}

func put(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.PUT(relativePath, hf(handlerFunc))
}

func patch(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.PATCH(relativePath, hf(handlerFunc))
}

func delete(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.DELETE(relativePath, hf(handlerFunc))
}

func options(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.OPTIONS(relativePath, hf(handlerFunc))
}

func head(group *gin.RouterGroup, relativePath string, handlerFunc handlerFunc) {
	group.HEAD(relativePath, hf(handlerFunc))
}

func hf(handlerFunc handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.New()

		err := handlerFunc(ctx, c)

		c.Set("request-id", ctx.RequestId())

		if err != nil {
			switch v := err.(type) {
			case *xerrors.Expected:
				if v.StatusOk() {
					return
				} else {
					c.JSON(v.StatusCode(), v.Message())
				}
			case *xerrors.Validation:
				c.JSON(http.StatusBadRequest, v)
			default:
				if gin.Mode() == gin.DebugMode {
					c.JSONP(http.StatusInternalServerError, map[string]string{"request_id": ctx.RequestId(), "error": v.Error()})
				} else {
					c.JSONP(http.StatusInternalServerError, map[string]string{"request_id": ctx.RequestId()})
				}
			}

			_ = c.Error(errors.Errorf("%+v", err))
		}
	}
}
