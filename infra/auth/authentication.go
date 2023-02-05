package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"auth-test/models"
	"auth-test/services"
)

func NewTokenAuthentication(secret string) TokenAuthentication {
	return TokenAuthentication{
		secret: secret,
	}
}

type TokenAuthentication struct {
	secret string
}

func (a TokenAuthentication) Sign(accessToken models.AccessTokenInput) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["sub"] = accessToken.AccountID()
	claims["email"] = accessToken.Email()
	claims["iat"] = accessToken.Now().Unix()

	exp := accessToken.ExpiredAt()
	claims["exp"] = time.Date(
		exp.Year(), exp.Month(), exp.Day(), exp.Hour(), exp.Minute(), exp.Second(), 0, exp.Location(),
	).Unix()
	signedToken, err := jwtToken.SignedString([]byte(a.secret))
	if err != nil {
		return "", services.NewApplicationErr(services.FailedSingedToken, err)
	}

	return signedToken, nil
}

func (a TokenAuthentication) Verify(token string) error {
	signedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, services.NewApplicationErr(services.InvalidToken, errors.New("証明の検証に失敗しました"))
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return err
	}

	if signedToken == nil {
		return services.NewApplicationErr(services.InvalidToken, errors.New("トークンが存在しません"))
	}

	claims, ok := signedToken.Claims.(jwt.MapClaims)
	if !ok {
		return services.NewApplicationErr(services.InvalidClaim, errors.New("クレームのキャストに失敗"))
	}

	if _, ok = claims["sub"].(string); !ok {
		return services.NewApplicationErr(services.InvalidIssued, errors.New("発行者のキャストに失敗"))
	}

	now := time.Now()
	ok = claims.VerifyExpiresAt(now.Unix(), false)
	if !ok {
		return services.NewApplicationErr(
			services.ExpiredToken, fmt.Errorf("有効期限: %s, 現在時刻: %d", claims["exp"], now.Unix()))
	}

	return nil
}
