package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/entities"
)

// var ErrNotValidData = errors.New("Not valid data")

func setBindError(c *gin.Context, err error) {
	c.Error(err).SetType(gin.ErrorTypeBind)
}

func setAnyError(c *gin.Context, err error) {
	c.Error(err).SetType(gin.ErrorTypeAny)
}

func abortWithStatusMSG(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{
		"msg": msg,
	})
}

// TODO: logger, kind of bind errors
func errorHandler(c *gin.Context) {
	c.Next()

	for _, ginErr := range c.Errors {
		if ginErr.IsType(gin.ErrorTypeBind) {
			log.Print("ErrNotValidData")
			abortWithStatusMSG(c, http.StatusBadRequest, ginErr.Err.Error())
			return
		}

		if ginErr.IsType(gin.ErrorTypeAny) {
			switch {
			case errors.Is(ginErr.Err, entities.ErrorUserHasAlreadyExist):
				log.Print("ErrorUserHasAlreadyExist")
				abortWithStatusMSG(c, http.StatusBadRequest, entities.ErrorUserHasAlreadyExist.Error())
				return
			}
		}

		log.Print("something went wrong")
		c.AbortWithStatus(http.StatusTeapot)
		return
	}
}
