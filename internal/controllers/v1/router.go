package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/v1adhope/time-tracker/docs"
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

// @title time-tracker API
// @version 1.0

// @host localhost:8081
// @BasePath /v1
func Handle(router *Router) {
	router.Handler.Use(gin.Recovery())

	router.Handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Log.Info("swagger available by this address http://localhost:8081/swagger/index.html")

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
