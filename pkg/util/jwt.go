package util

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret = []byte("123456")
)

type Claims struct {
	Userid int64 `json:"userid"`
	Role   byte  `json:role, omitempty`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(userid int64) (string, error) {
	expiredAt := time.Now().Add(7 * 24 * time.Hour)

	claims := Claims{
		Userid: 1,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
