package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"auth-test/models"
	"auth-test/services"
)

type UserAccount struct {
	Service services.UserAccount
}

type InputUserAccount struct {
	Email    string `json:"email" required:"true"`
	Name     string `json:"name" required:"true"`
	Password string `json:"password" required:"true"`
}

func (a UserAccount) Get(c *gin.Context) {
	accountID := c.Param("id")
	userAccount, err := a.Service.Find(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("%v", userAccount)
	c.JSON(http.StatusOK, userAccount)
}

func (a UserAccount) List(c *gin.Context) {
	accounts, err := a.Service.List()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (a *UserAccount) Create(c *gin.Context) {
	var account InputUserAccount
	err := c.Bind(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := a.Service.Create(models.UserAccount{
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

func (a *UserAccount) Update(c *gin.Context) {
	id := c.Param("id")
	var account InputUserAccount
	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	result, err := a.Service.Update(models.UserAccount{
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

func (a UserAccount) Delete(c *gin.Context) {
	id := c.Param("id")
	err := a.Service.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
