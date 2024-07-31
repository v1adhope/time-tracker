package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/pkg/logger"
)

func trackingHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Track(c.Request.URL.RequestURI(), c.ClientIP(), c.Writer.Status())
	}
}
