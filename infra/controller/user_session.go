package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"auth-test/services"
)

func NewJWTAuth(service services.Authorizer) JWTAuth {
	return JWTAuth{
		authorizer: service,
	}
}

type loginForm struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type JWTAuth struct {
	authorizer services.Authorizer
}

func (s JWTAuth) Login(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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

func (s JWTAuth) CheckAuthentication(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithStatus(http.StatusBadRequest)
		c.Writer.Write([]byte("empty token"))
		return
	}

	err := s.authorizer.Verify(token)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		c.Writer.Write([]byte(err.Error()))
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
		c.JSON(http.StatusInternalServerError, err.Error())
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

	id := c.Param("id")
	err := s.authorizer.FindUser(id, token)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Next()
}

func (s StoredAuth) Logout(c *gin.Context) {
	t := c.GetHeader("Authorization")
	token := strings.Replace(t, "Bearer ", "", 1)
	if err := s.authorizer.SignOut(token); err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	return
}
