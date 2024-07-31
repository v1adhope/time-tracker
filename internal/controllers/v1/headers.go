package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	HeaderLocation = "Location"
)

func parseBaseReqURL(c *gin.Context) string {
	scheme := "http"

	if c.Request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

func setLocationHeader(c *gin.Context, location string) {
	c.Header(HeaderLocation, location)
}

func setUserLocationHeader(c *gin.Context, id string) {
	location := fmt.Sprintf(
		"%s/v1/users?id=%s",
		parseBaseReqURL(c),
		id,
	)

	setLocationHeader(c, location)
}

func setTaskLocationHeader(c *gin.Context, id string) {
	location := fmt.Sprintf(
		"%s/v1/tasks/summary-time/%s",
		parseBaseReqURL(c),
		id,
	)

	setLocationHeader(c, location)
}
