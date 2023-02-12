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

func NewSessionAuth(service services.UserSession) UserSessionHandler {
	return UserSessionHandler{
		session: service,
	}
}

type loginForm struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,nist_sp_800_63,min=8,max=72" example:"string"`
}

type UserSessionHandler struct {
	session services.UserSession
}

type SessionToken struct {
	Value string `json:"value" binding:"required,uuid" example:"12345678-89ab-cdef-ghij-klmopqrstuvw"`
}

// Login get session token
// @Summary Return session token for login user
// @Tags Login
// @Param loginFrom body controller.loginForm true "Email and Password"
// @Produce json
// @Success 200 {object} controller.SessionToken
// @Failure default {object} controller.errResponse
// @Router  /session/login [post]
func (a UserSessionHandler) Login(c *gin.Context) {
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	token, err := a.session.Sign(form.Email, form.Password, uuid.New().String(), time.Now())
	if err != nil {
		status, response := newErrResponse(err, form.Email)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.JSON(http.StatusOK, SessionToken{Value: token})
	return
}

func (a UserSessionHandler) CheckAuthenticatedOwner(c *gin.Context) {
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			errResponse{Message: services.EmptyToken.Error(), Detail: "トークンは必須です"},
		)
		return
	}

	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	err := a.session.FindOwner(params.ID, token)
	if err != nil {
		status, response := newErrResponse(err, token)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.Next()
}

// Logout delete session token
// @Summary Return status by delete session
// @Tags Logout
// @securityDefinitions.apiKey ApiKeyAuth
// @Param id path string true "User ID by UUID"
// @Produce json
// @Success 200
// @Failure 400 {object} controller.errResponse
// @Failure 401 {object} controller.errResponse
// @Failure 500 {object} controller.errResponse
// @Router  /session/logout/{id} [delete]
// @Security Bearer
func (a UserSessionHandler) Logout(c *gin.Context) {
	t := c.GetHeader("Authorization")
	token := strings.Replace(t, "Bearer ", "", 1)
	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	if err := a.session.SignOut(params.ID, token); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	return
}
