package middleware

import (
	"errors"
	"net/http"
	"strings"

	"authkey/pkg/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	ErrTokenEmpty   = errors.New("empty token")
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenTimeout = errors.New("token timeout")

	_token = "tk"
)

func GetJWT(c *gin.Context) *util.Claims {
	return c.MustGet(_token).(*util.Claims)
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})

			c.Abort()
			return
		}

		c.Set(_token, token)

		c.Next()
	}
}
