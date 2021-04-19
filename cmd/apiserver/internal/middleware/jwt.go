package middleware

import (
	"errors"
	"net/http"
	"strings"

	"authkey/pkg/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/meilihao/water"
)

var (
	ErrTokenEmpty   = errors.New("empty token")
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenTimeout = errors.New("token timeout")

	_token = "tk"
)

func GetJWT(c *water.Context) *util.Claims {
	return c.Environ.Get(_token).(*util.Claims)
}

func JWT() water.HandlerFunc {
	return func(c *water.Context) {
		raw := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")
		if raw == "" {
			raw = c.Query("token")
		}

		var err error
		var token *util.Claims

		if raw == "" {
			err = util.I18nError(c, "token.empty")
		} else {
			token, err = util.ParseToken(raw)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = util.I18nError(c, "token.expired")
				default:
					err = util.I18nError(c, "token.invalid")
				}
			}
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, water.H{
				"error": err,
			})

			return
		}

		c.Environ.Set(_token, token)

		c.Next()
	}
}
