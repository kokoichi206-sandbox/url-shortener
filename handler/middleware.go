package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kokoichi206-sandbox/url-shortener/util"
)

func (h *handler) requestIDMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := util.WithRequestID(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
