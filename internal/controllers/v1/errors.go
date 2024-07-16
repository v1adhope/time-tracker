package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/entities"
)

// var ErrNotValidData = errors.New("Not valid data")
// var ErrNoID = errors.New("missing id")

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

// TODO: logger, kinds of bind errors
func errorHandler(c *gin.Context) {
	c.Next()

	for _, ginErr := range c.Errors {
		switch ginErr.Type {
		case gin.ErrorTypeBind:
			log.Print("ErrNotValidData")
			abortWithStatusMSG(c, http.StatusBadRequest, ginErr.Err.Error())
			return
		case gin.ErrorTypeAny:
			switch {
			case errors.Is(ginErr.Err, entities.ErrorUserHasAlreadyExist):
				log.Print("ErrorUserHasAlreadyExist")
				abortWithStatusMSG(c, http.StatusBadRequest, entities.ErrorUserHasAlreadyExist.Error())
				return
			case errors.Is(ginErr.Err, entities.ErrorUserDoesNotExist),
				errors.Is(ginErr.Err, entities.ErrorTaskDoesNotExist),
				errors.Is(ginErr.Err, entities.ErrorNoAnyTasksForThisUser):

				log.Print("ErrorUserDoesNotExist or ErrorTaskDoesNotExist or ErrorNoAnyTasksForThisUser")
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		log.Printf("something went wrong: %s", ginErr.Err.Error())
		c.AbortWithStatus(http.StatusTeapot)
		return
	}
}
