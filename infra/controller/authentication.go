package controller

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"auth-test/services"
)

type RefreshToken struct {
	Value string `json:"value" binding:"required,uuid"`
}

type OauthToken struct {
	Access  string `json:"access" binding:"required"`
	Refresh string `json:"refresh" binding:"required,uuid"`
}

func NewTokenHandler(service services.OauthToken) TokenHandler {
	return TokenHandler{
		authenticateSvc: service,
	}
}

type TokenHandler struct {
	authenticateSvc services.OauthToken
}

func (h TokenHandler) Claim(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	token, err := h.authenticateSvc.Claim(form.Email, form.Password, uuid.New().String(), time.Now())
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, OauthToken{Access: token.Access(), Refresh: token.Refresh()})
}

func (h TokenHandler) Refresh(c *gin.Context) {
	var refresh RefreshToken
	err := c.Bind(&refresh)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.JSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	now := time.Now()
	token, err := h.authenticateSvc.Refresh(uuid.New().String(), refresh.Value, now.UTC())
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}

	c.JSON(http.StatusOK, OauthToken{Access: token.Access(), Refresh: token.Refresh()})
}

func (h TokenHandler) VerifyAccessToken(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty token"))
		return
	}

	err := h.authenticateSvc.Verify(token)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}

	c.Next()
}
