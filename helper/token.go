package helper

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtStruct struct {
	Id       string
	Username string
	Role     string
	jwt.RegisteredClaims
}

func NewToken(g *JwtStruct, secretKey string, expire int64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		JwtStruct{
			Id:       g.Id,
			Username: g.Username,
			Role:     g.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire * 1000 * 1000))),
			},
		},
	)

	return claims.SignedString([]byte(secretKey))
}
