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
		users.DELETE("/:id", router.Delete)
		users.PATCH("/:id", router.Update)
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

type deleteUserReqParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (r *UserRouter) Delete(c *gin.Context) {
	params := deleteUserReqParams{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	if err := r.userUsecase.Delete(c.Request.Context(), params.ID); err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

type updateUserReqParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type updateUserReq struct {
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	PassportNumber string `json:"passportNumber" binding:"len=9"`
}

func (r *UserRouter) Update(c *gin.Context) {
	params := updateUserReqParams{}

	if err := c.ShouldBindUri(&params); err != nil {
		setBindError(c, err)
		return
	}

	req := updateUserReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		setBindError(c, err)
		return
	}

	if err := r.userUsecase.Update(c.Request.Context(), entities.User{
		ID:             params.ID,
		Surname:        req.Surname,
		Name:           req.Name,
		Patronymic:     req.Patronymic,
		Address:        req.Address,
		PassportNumber: req.PassportNumber,
	}); err != nil {
		setAnyError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
