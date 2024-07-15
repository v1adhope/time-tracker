package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/usecases"
)

type Router struct {
	Handler  *gin.Engine
	Usecases *usecases.Usecases
}

func Handle(router *Router) {
	router.Handler.Use(gin.Logger(), gin.Recovery())

	v1 := router.Handler.Group("/v1")

	v1.Use(errorHandler)
	{
		handleUser(&UserRouter{
			handler:     v1,
			userUsecase: router.Usecases.User,
		})
	}
}
