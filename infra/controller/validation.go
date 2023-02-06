package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/s-tajima/nspv"
)

const (
	uuidTokenFormat = "フォーマットを 12345678-89ab-cdef-ghij-klmopqrstuvw にして下さい"
)

var (
	passwordValidate = nspv.NewValidator()
)

func newPathParamError(err validator.FieldError) pathParamError {
	return pathParamError{
		err: err,
	}
}

type pathParamError struct {
	err validator.FieldError
}

func (e pathParamError) getResponse() errResponse {
	var response errResponse
	errorMsg := e.err
	switch errorMsg.Field() {
	case "ID":
		response = newValidationErr(
			fmt.Sprintf("ユーザIDの値 %s は不正なフォーマットです", errorMsg.Value()),
			uuidTokenFormat,
		)
	}

	return response
}

func ValidatePassword(fl validator.FieldLevel) bool {
	passwd := fl.Field().String()
	result, err := passwordValidate.Validate(passwd)
	if err != nil {
		return false
	}

	if result != nspv.Ok {
		return false
	}

	return true
}

func newAccountBodyError(err validator.FieldError) AuthBodyError {
	return AuthBodyError{
		err: err,
	}
}

type AuthBodyError struct {
	err validator.FieldError
}

// TODO バリデーションエラーの詳細な情報を情報を取得する方法を調査する
func (e AuthBodyError) getResponse() errResponse {
	var response errResponse
	errorMsg := e.err
	switch errorMsg.Field() {
	case "Password":
		response = newValidationErr(
			"パスワードが不正なフォーマットです",
			"",
		)
	case "Email":
		response = newValidationErr(
			fmt.Sprintf("emailアドレス %s は不正なフォーマットです", errorMsg.Value()),
			"",
		)
	case "Value":
		response = newValidationErr(
			fmt.Sprintf("リフレッシュトークン %s は不正なフォーマットです", errorMsg.Value()),
			uuidTokenFormat,
		)
	}

	return response
}
