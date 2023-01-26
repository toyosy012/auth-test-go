package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"auth-test/models"
	"auth-test/services"
)

type UserAccountHandler struct {
	UserAccountService services.UserAccount
}

type InputUserAccount struct {
	Email    string `json:"email" required:"true"`
	Name     string `json:"name" required:"true"`
	Password string `json:"password" required:"true"`
}

func (h UserAccountHandler) Get(c *gin.Context) {
	accountID := c.Param("id")
	userAccount, err := h.UserAccountService.Find(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, userAccount)
}

func (h UserAccountHandler) List(c *gin.Context) {
	accounts, err := h.UserAccountService.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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

	result, err := h.UserAccountService.Create(models.UserAccount{
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
	id := c.Param("id")
	var account InputUserAccount
	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.UserAccountService.Update(models.UserAccount{
		ID:       id,
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
	id := c.Param("id")
	err := h.UserAccountService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
