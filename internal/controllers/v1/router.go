package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/usecases"
	"github.com/v1adhope/time-tracker/pkg/logger"
)

type Config struct {
	Mode string `koanf:"APP_GIN_MODE"`
}

type Router struct {
	Handler  *gin.Engine
	Usecases *usecases.Usecases
	Log      logger.Logger
}

func Handle(router *Router) {
	router.Handler.Use(gin.Recovery())

	v1 := router.Handler.Group("/v1")

	v1.Use(errorHandler(router.Log))
	{
		handleUser(&userRouter{
			handler:     v1,
			userUsecase: router.Usecases.User,
		})
		handleTask(&taskRouter{
			handler:     v1,
			taskUsecase: router.Usecases.Task,
		})
	}
}
