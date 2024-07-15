package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/entities"
	"github.com/v1adhope/time-tracker/internal/usecases"
)

type UserRouter struct {
	handler     *gin.RouterGroup
	userUsecase usecases.User
}

func handleUser(router *UserRouter) {
	users := router.handler.Group("/users")
	{
		users.POST("/", router.Create)
	}
}

type createUserReq struct {
	Surname        string `json:"surname" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Patronymic     string `json:"patronymic" binding:"required"`
	Address        string `json:"address" binding:"required"`
	PassportNumber string `json:"passportNumber" binding:"required,len=9"`
}

func (r *UserRouter) Create(c *gin.Context) {
	req := createUserReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	user, err := r.userUsecase.Create(c.Request.Context(), entities.User{
		Surname:        req.Surname,
		Name:           req.Name,
		Patronymic:     req.Patronymic,
		Address:        req.Address,
		PassportNumber: req.PassportNumber,
	})
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}
