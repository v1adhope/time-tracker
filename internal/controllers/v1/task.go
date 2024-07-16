package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/usecases"
)

type taskRouter struct {
	handler     *gin.RouterGroup
	taskUsecase usecases.Task
}

func handleTask(router *taskRouter) {
	tasks := router.handler.Group("/tasks")
	{
		tasks.POST("/start/:id", router.Start)
	}
}

type startTaskReqParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (r *taskRouter) Start(c *gin.Context) {
	params := startTaskReqParams{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	task, err := r.taskUsecase.Start(c.Request.Context(), params.ID)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}
