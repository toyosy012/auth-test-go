package services

import (
	"errors"
)

// infra層で利用するエラーの詳細を伝えるメッセージ変数
var (
	NoTokenRecord      = errors.New("トークンデータは存在しません")
	DuplicateToken     = errors.New("トークンが既に存在します")
	NoUsersRecord      = errors.New("ユーザリストの取得に失敗")
	NoUserRecord       = errors.New("ユーザ情報が取得できません")
	NoUserEmail        = errors.New("emailで登録されたユーザは存在しません")
	DuplicateUserEmail = errors.New("emailが既に存在しています")
	EmptyToken         = errors.New("トークンが存在しません")
	ExpiredToken       = errors.New("有効期限切れトークンです")
	NoSessionRecord    = errors.New("セッションは存在しません")
	InvalidToken       = errors.New("無効なトークンです")
	InvalidClaim       = errors.New("無効なペイロードです")
	FailedSingedToken  = errors.New("トークンの署名に失敗しました")
	InvalidIssued      = errors.New("発行時期が無効なトークンです")
	TooLongPassword    = errors.New("パスワードを72文字以内にしてください")
	InvalidUUIDFormat  = errors.New("無効なUUIDです")
	InternalServerErr  = errors.New("サーバエラーが発生しました")
)

// services層で利用するエラーが発生したユースケースを伝えるメッセージ変数
var (
	FailedShowUser     = errors.New("ユーザの取得に失敗")
	FailedListUser     = errors.New("ユーザリストの取得に失敗")
	FailedCreateUser   = errors.New("ユーザ登録に失敗")
	FailedUpdateUser   = errors.New("ユーザ情報の更新に失敗")
	FailedDeleteUser   = errors.New("ユーザ削除に失敗")
	FailedAuthenticate = errors.New("認証に失敗しました")
	FailedCreateToken  = errors.New("トークン作成に失敗しました")
	FailedCheckLogin   = errors.New("ログイン情報が確認できませんでした")
	FailedLogin        = errors.New("ログインに失敗しました")
	FailedLogout       = errors.New("ログアウトに失敗しました")
)

func NewApplicationErr(message, detail error) ApplicationErr {
	return ApplicationErr{Message: message, Detail: detail}
}

type ApplicationErr struct {
	Message error
	Detail  error
}

func (e ApplicationErr) Error() string { return e.Message.Error() }
func (e ApplicationErr) Unwrap() error { return e.Detail }
