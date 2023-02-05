package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

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

type userPathParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type inputUserAccount struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,nist_sp_800_63"`
}

type userAccountResponse struct {
	ID    string `json:"id" binding:"required,uuid"`
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (h UserAccountHandler) Get(c *gin.Context) {
	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}

	userAccount, err := h.service.Find(params.ID)
	if err != nil {
		status, response := newErrResponse(err, params.ID)
		c.AbortWithStatusJSON(status, response)
	}

	c.JSON(http.StatusOK, userAccountResponse{ID: userAccount.ID(), Email: userAccount.Email(), Name: userAccount.Name()})
}

func (h UserAccountHandler) List(c *gin.Context) {
	accounts, err := h.service.List()
	if err != nil {
		status, response := newErrResponse(err, "")
		c.AbortWithStatusJSON(status, response)
		return
	}

	var results []userAccountResponse
	for _, a := range accounts {
		results = append(results, userAccountResponse{ID: a.ID(), Email: a.Email(), Name: a.Name()})
	}

	c.JSON(http.StatusOK, results)
}

func (h *UserAccountHandler) Create(c *gin.Context) {
	var account inputUserAccount
	err := c.Bind(&account)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	result, err := h.service.Create(
		models.NewUserAccount(uuid.New().String(), account.Email, account.Name, account.Password),
	)
	if err != nil {
		status, response := newErrResponse(err, account.Email)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.JSON(http.StatusOK, userAccountResponse{ID: result.ID(), Email: result.Email(), Name: result.Name()})
}

func (h *UserAccountHandler) Update(c *gin.Context) {
	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	var account inputUserAccount
	err := c.BindJSON(&account)
	if err != nil {
		accountBodyParam := newAccountBodyError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, accountBodyParam.getResponse())
		return
	}

	result, err := h.service.Update(
		models.NewUserAccount(params.ID, account.Email, account.Name, account.Password),
	)

	if err != nil {
		status, response := newErrResponse(err, account.Email)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.JSON(http.StatusOK, userAccountResponse{ID: result.ID(), Email: result.Email(), Name: result.Name()})
}

func (h UserAccountHandler) Delete(c *gin.Context) {
	var params userPathParams
	if err := c.BindUri(&params); err != nil {
		pathParamErr := newPathParamError(err.(validator.ValidationErrors)[0])
		c.AbortWithStatusJSON(http.StatusBadRequest, pathParamErr.getResponse())
		return
	}
	err := h.service.Delete(params.ID)
	if err != nil {
		status, response := newErrResponse(err, params.ID)
		c.AbortWithStatusJSON(status, response)
		return
	}

	c.Status(http.StatusOK)
}
