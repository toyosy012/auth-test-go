package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"auth-test/services"
)

func NewLogin(service services.Authorizer) LoginHandler {
	return LoginHandler{
		authorizer: service,
	}
}

type loginForm struct {
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type LoginHandler struct {
	authorizer services.Authorizer
}

func (h LoginHandler) Login(c *gin.Context) {
	now := time.Now()
	var form loginForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.authorizer.Sign(form.Email, form.Password, now)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.Status(http.StatusOK)
	return
}

func (h LoginHandler) CheckAuthentication(c *gin.Context) {
	now := time.Now()
	t := c.GetHeader("Authorization")

	token := strings.Replace(t, "Bearer ", "", 1)
	if "" == token {
		c.AbortWithStatus(http.StatusBadRequest)
		c.Writer.Write([]byte("empty token"))
		return
	}

	err := h.authorizer.Verify(token, now)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	c.Next()
}
