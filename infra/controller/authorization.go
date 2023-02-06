package controller

import (
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

type AuthToken struct {
	IDToken string `json:"id" binding:"required"`
	Refresh string `json:"refresh" binding:"required,uuid"`
}

func NewTokenHandler(service services.Authorizer) TokenHandler {
	return TokenHandler{
		authenticateSvc: service,
	}
}

type TokenHandler struct {
	authenticateSvc services.Authorizer
}

func (h TokenHandler) Claim(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	token, err := h.authenticateSvc.Claim(form.Email, form.Password, uuid.New().String(), time.Now())
	if err != nil {
		status, response := newErrResponse(err, form.Email)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.JSON(http.StatusOK, AuthToken{IDToken: token.IDToken(), Refresh: token.Refresh()})
}

func (h TokenHandler) Refresh(c *gin.Context) {
	var refresh RefreshToken
	err := c.Bind(&refresh)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	now := time.Now()
	token, err := h.authenticateSvc.Refresh(uuid.New().String(), refresh.Value, now.UTC())
	if err != nil {
		status, response := newErrResponse(err, refresh.Value)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.JSON(http.StatusOK, AuthToken{IDToken: token.IDToken(), Refresh: token.Refresh()})
}

func (h TokenHandler) VerifyIDToken(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errResponse{Message: services.EmptyToken.Error(), Detail: "トークンは必須です"},
		)
		return
	}

	err := h.authenticateSvc.Verify(token)
	if err != nil {
		status, response := newErrResponse(err, "")
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.Next()
}
