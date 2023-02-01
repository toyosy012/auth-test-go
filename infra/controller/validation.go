package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type validateError interface {
	getResponse() errorResponse
}

func newPathParamError(err validator.FieldError) pathParamError {
	return pathParamError{
		err: err,
	}
}

type pathParamError struct {
	err validator.FieldError
}

func (e pathParamError) getResponse() errorResponse {
	var response errorResponse
	switch e.err.Field() {
	case "ID":
		response = newErrorResponse(
			fmt.Sprintf("ユーザIDの値 %s は不正なフォーマットです", e.err.Value()),
			"フォーマットを 12345678-89ab-cdef-ghij-klmopqrstuvw にして下さい",
		)
	}

	return response
}
