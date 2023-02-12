package controller

import (
	"errors"
	"fmt"
	"net/http"

	"auth-test/services"
)

func newValidationErr(message string, detail string) errResponse {
	return errResponse{
		Message: message,
		Detail:  detail,
	}
}

type errResponse struct {
	Message string `json:"message" binding:"required"`
	Detail  string `json:"detail,omitempty"`
}

func (e errResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Detail)
}

func newErrResponse(err error, detail string) (status int, responseErr errResponse) {
	applicationErr := errors.Unwrap(err)
	switch {
	case errors.Is(applicationErr, services.InternalServerErr):
		status = http.StatusInternalServerError
	case errors.Is(applicationErr, services.NoSessionRecord):
		status = http.StatusUnauthorized
	default:
		status = http.StatusBadRequest
	}

	var detailMsg string
	if detail != "" {
		detailMsg = fmt.Sprintf(": %s", detail)
	}

	return status, errResponse{
		Message: err.Error(),
		Detail:  fmt.Sprintf("%s%s", errors.Unwrap(err).Error(), detailMsg),
	}
}
