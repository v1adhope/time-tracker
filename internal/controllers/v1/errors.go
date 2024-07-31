package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/entities"
	"github.com/v1adhope/time-tracker/pkg/logger"
)

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

func errorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			switch ginErr.Type {
			case gin.ErrorTypeBind:
				log.Debug(ginErr.Err)
				abortWithStatusMSG(c, http.StatusBadRequest, ginErr.Err.Error())
				return
			case gin.ErrorTypeAny:
				switch {
				case errors.Is(ginErr.Err, entities.ErrorUserHasAlreadyExistWithThatPassport):
					log.Debug(ginErr.Err)
					abortWithStatusMSG(c, http.StatusBadRequest, ginErr.Err.Error())
					return
				case errors.Is(ginErr.Err, entities.ErrorUsersDoesNotExist),
					errors.Is(ginErr.Err, entities.ErrorTaskDoesNotExist),
					errors.Is(ginErr.Err, entities.ErrorNoAnyTasksForThisUser):

					log.Debug(ginErr.Err)
					c.AbortWithStatus(http.StatusNoContent)
					return
				case errors.Is(ginErr.Err, entities.ErrorUserDoesNotExistWithThatPassportInfoExeption):
					log.Debug(ginErr.Err)
					abortWithStatusMSG(c, http.StatusBadRequest, ginErr.Err.Error())
					return
				}
			}

			log.Error(ginErr.Err)
			c.AbortWithStatus(http.StatusTeapot)
			return
		}
	}
}
