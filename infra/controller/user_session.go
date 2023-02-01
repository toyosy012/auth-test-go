package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"auth-test/services"
)

func NewJWTAuth(service services.Authorizer) JWTAuth {
	return JWTAuth{
		authorizer: service,
	}
}

type loginForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,nist_sp_800_63"`
}

type JWTAuth struct {
	authorizer services.Authorizer
}

func (s JWTAuth) Login(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	token, err := s.authorizer.Sign(form.Email, form.Password)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
	return
}

func (s JWTAuth) CheckAuthentication(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty token"))
		return
	}

	err := s.authorizer.Verify(token)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}

	c.Next()
}

func NewStoredAuth(s services.StoredAuthorizer) StoredAuth {
	return StoredAuth{
		authorizer: s,
	}
}

type StoredAuth struct {
	authorizer services.StoredAuthorizer
}

func (s StoredAuth) Login(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	token, err := s.authorizer.Sign(form.Email, form.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
	return
}

func (s StoredAuth) CheckAuthentication(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty token"))
		return
	}

	err := s.authorizer.Verify(token)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Next()
}

func (s StoredAuth) CheckAuthenticatedOwner(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty token"))
		return
	}

	var params UserPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	err := s.authorizer.FindUser(params.ID, token)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Next()
}

func (s StoredAuth) Logout(c *gin.Context) {
	t := c.GetHeader("Authorization")
	token := strings.Replace(t, "Bearer ", "", 1)
	var params UserPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	if err := s.authorizer.SignOut(params.ID, token); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	return
}
