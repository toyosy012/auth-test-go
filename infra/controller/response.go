package controller

import (
	"fmt"
)

func newErrorResponse(message string, detail string) errorResponse {
	return errorResponse{
		Message: message,
		Detail:  detail,
	}
}

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (r errorResponse) Error() string {
	return fmt.Sprintf("%s: %s", r.Message, r.Detail)
}
