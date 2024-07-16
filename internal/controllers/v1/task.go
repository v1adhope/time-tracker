package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/entities"
	"github.com/v1adhope/time-tracker/internal/usecases"
)

type taskRouter struct {
	handler     *gin.RouterGroup
	taskUsecase usecases.Task
}

func handleTask(router *taskRouter) {
	tasks := router.handler.Group("/tasks")
	{
		tasks.POST("/start/:userId", router.Start)
		tasks.PATCH("/end/:id", router.End)
		tasks.GET("/summary-time/:userId", router.SummaryTime)
	}
}

type startTaskReqParams struct {
	UserID string `uri:"userId" binding:"required,uuid"`
}

func (r *taskRouter) Start(c *gin.Context) {
	params := startTaskReqParams{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	task, err := r.taskUsecase.Start(c.Request.Context(), params.UserID)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

type endTaskReqParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (r *taskRouter) End(c *gin.Context) {
	params := endTaskReqParams{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	finishedAt, err := r.taskUsecase.End(c.Request.Context(), params.ID)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"finishedAt": finishedAt,
	})
}

type summaryTimeReqParams struct {
	UserID string `uri:"userId" binding:"required,uuid"`
}

type summaryTimeReqQuery struct {
	StartTime string `form:"startTime" binding:"omitempty,sorttime"`
	EndTime   string `form:"endTime" binding:"omitempty,sorttime"`
}

func (r *taskRouter) SummaryTime(c *gin.Context) {
	params := summaryTimeReqParams{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	query := summaryTimeReqQuery{}

	if err := c.ShouldBindQuery(&query); err != nil {
		setBindError(c, err)
		return
	}

	tasks, err := r.taskUsecase.GetReportSummaryTime(
		c.Request.Context(),
		params.UserID,
		entities.TaskSort{
			StartTime: query.StartTime,
			EndTime:   query.EndTime,
		},
	)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}
