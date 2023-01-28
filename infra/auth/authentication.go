package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"auth-test/models"
)

func NewTokenAuthentication(secret string, availabilityTime time.Duration) TokenAuthentication {
	return TokenAuthentication{
		secret:           secret,
		availabilityTime: availabilityTime,
	}
}

type TokenAuthentication struct {
	secret           string
	availabilityTime time.Duration
}

func (a TokenAuthentication) Sign(account models.UserAccount) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = account.ID
	claims["email"] = account.Email
	now := time.Now()
	claims["iat"] = now.Unix()

	exp := now.Add(a.availabilityTime)
	claims["exp"] = time.Date(exp.Year(), exp.Month(), exp.Day(), exp.Hour(), exp.Minute(), exp.Second(), 0, exp.Location()).Unix()
	signedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", errors.New("")
	}

	return signedToken, nil
}

func (a TokenAuthentication) Verify(token string) error {
	signedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]))
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return err
	}

	if signedToken == nil {
		return errors.New("bad request token")
	}

	claims, ok := signedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("bad request token")
	}

	if _, ok := claims["sub"].(string); !ok {
		return errors.New("not found user")
	}

	now := time.Now()
	ok = claims.VerifyExpiresAt(now.Unix(), false)
	if !ok {
		return errors.New("authentication expired")
	}

	return nil
}
