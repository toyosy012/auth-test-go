package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"auth-test/services"
)

func NewSessionAuth(service services.UserSession) UserSession {
	return UserSession{
		session: service,
	}
}

type loginForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,nist_sp_800_63"`
}

type UserSession struct {
	session services.UserSession
}

func (a UserSession) Login(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	token, err := a.session.Sign(form.Email, form.Password, uuid.New().String(), time.Now())
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
	return
}

func (a UserSession) CheckAuthenticatedOwner(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty token"))
		return
	}

	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	err := a.session.FindOwner(params.ID, token)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Next()
}

func (a UserSession) Logout(c *gin.Context) {
	t := c.GetHeader("Authorization")
	token := strings.Replace(t, "Bearer ", "", 1)
	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	if err := a.session.SignOut(params.ID, token); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	return
}
