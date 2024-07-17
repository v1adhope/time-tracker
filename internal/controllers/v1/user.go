package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/entities"
	"github.com/v1adhope/time-tracker/internal/usecases"
)

type userRouter struct {
	handler     *gin.RouterGroup
	userUsecase usecases.User
}

func handleUser(router *userRouter) {
	users := router.handler.Group("/users")
	{
		users.POST("/", router.Create)
		users.DELETE("/:id", router.Delete)
		users.PATCH("/:id", router.Update)
		users.GET("/", router.All)
		users.GET("/info", router.Info)
	}
}

type createUserReq struct {
	Surname        string `json:"surname" binding:"required" example:"Bode"`
	Name           string `json:"name" binding:"required" example:"Rogers"`
	Patronymic     string `json:"patronymic" binding:"required" example:"Robertovich"`
	Address        string `json:"address" binding:"required" example:"1123 Ola Brook"`
	PassportNumber string `json:"passportNumber" binding:"required,passport" example:"6666 666666"`
}

// @tags users
// @summary Create user
// @accept json
// @param user body createUserReq true "User request model"
// @response 201 {object} entities.User
// @response 400
// @response 500
// @router /users [post]
func (r *userRouter) Create(c *gin.Context) {
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

// @tags users
// @summary Delete user
// @param id path string true "User id (uuid)"
// @response 200
// @response 204 "There's no user to delete"
// @response 400
// @response 500
// @router /users/{id} [delete]
func (r *userRouter) Delete(c *gin.Context) {
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
	Surname        string `json:"surname" example:"Wyman"`
	Name           string `json:"name" example:"Nicholas"`
	Patronymic     string `json:"patronymic" example:"Victorovich"`
	Address        string `json:"address" example:"516 Carlee Statio"`
	PassportNumber string `json:"passportNumber" binding:"passport" example:"7777 777777"`
}

// @tags users
// @summary Update user
// @param id path string true "User id (uuid)"
// @accept json
// @param user body updateUserReq true "User request model"
// @response 200
// @response 204 "There's no user to change"
// @response 400
// @response 500
// @router /users/{id} [patch]
func (r *userRouter) Update(c *gin.Context) {
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

type allUserQuery struct {
	Surname        string `form:"surname" binding:"omitempty,filterstring"`
	Name           string `form:"name" binding:"omitempty,filterstring"`
	Patronymic     string `form:"patronymic" binding:"omitempty,filterstring"`
	Address        string `form:"address" binding:"omitempty,filterstring"`
	PassportNumber string `form:"passportNumber" binding:"omitempty,filterstring"`

	Limit  string `form:"limit" binding:"omitempty,number"`
	Offset string `form:"offset" binding:"omitempty,number"`
}

// @tags users
// @summary Get all users
// @param limit query uint64 false "Pagination control"
// @param offset query uint64 false "Pagination control"
// @param surname query string false "Custom type consitst operation:value. Allowed operations eq, ilike"
// @param name query string false "Custom type consitst operation:value. Allowed operations eq, ilike"
// @param address query string false "Custom type consitst operation:value. Allowed operations eq, ilike"
// @param passportNumber query string false "Custom type consitst operation:value. Allowed operations eq, ilike"
// @response 200
// @response 204 "No any users by this request"
// @response 400
// @response 500
// @router /users [get]
func (r *userRouter) All(c *gin.Context) {
	query := allUserQuery{}

	if err := c.ShouldBindQuery(&query); err != nil {
		setBindError(c, err)
		return
	}

	users, err := r.userUsecase.GetAll(c.Request.Context(), entities.UserRepresentation{
		Pagination: entities.UserPagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
		Filter: entities.UserFilter{
			BySurname:        query.Surname,
			ByName:           query.Name,
			ByPatronymic:     query.Patronymic,
			ByAddress:        query.Address,
			ByPassportNumber: query.PassportNumber,
		},
	})
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

type infoUserQuery struct {
	PassportSeries string `form:"passportSeries" binding:"required,len=4"`
	PassportNumber string `form:"passportNumber" binding:"required,len=6"`
}

// @tags users
// @summary Info endpoint
// @param passportSeries query string true "Should be number len=4"
// @param passportNumber query string true "Should be number len=6"
// @response 200 "Consist empty user"
// @response 400
// @response 500
// @router /users/info [get]
func (r *userRouter) Info(c *gin.Context) {
	query := infoUserQuery{}

	if err := c.ShouldBindQuery(&query); err != nil {
		setBindError(c, err)
		return
	}

	passportField := fmt.Sprintf("%s %s", query.PassportSeries, query.PassportNumber)

	user, err := r.userUsecase.Get(c.Request.Context(), passportField)
	if err != nil {
		setAnyError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
