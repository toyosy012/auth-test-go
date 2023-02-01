package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"auth-test/models"
	"auth-test/services"
)

func NewUserAccountHandler(svc services.UserAccount, validate validator.Validate) UserAccountHandler {
	return UserAccountHandler{
		service:  svc,
		validate: validate,
	}
}

type UserAccountHandler struct {
	service  services.UserAccount
	validate validator.Validate
}

type UserPathParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type InputUserAccount struct {
	Email    string `json:"email" required:"true"`
	Name     string `json:"name" required:"true"`
	Password string `json:"password" required:"true"`
}

func (h UserAccountHandler) Get(c *gin.Context) {
	var params UserPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	userAccount, err := h.service.Find(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, userAccount)
}

func (h UserAccountHandler) List(c *gin.Context) {
	accounts, err := h.service.List()
	if err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *UserAccountHandler) Create(c *gin.Context) {
	var account InputUserAccount
	err := c.Bind(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Create(models.UserAccount{
		Email:    account.Email,
		Name:     account.Name,
		Password: account.Password,
	})

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserAccountHandler) Update(c *gin.Context) {
	var params UserPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	var account InputUserAccount
	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Update(models.UserAccount{
		ID:       params.ID,
		Email:    account.Email,
		Name:     account.Name,
		Password: account.Password,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h UserAccountHandler) Delete(c *gin.Context) {
	var params UserPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	err := h.service.Delete(params.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
