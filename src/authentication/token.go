package authentication

import (
	"api/src/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(userId uint64) (string, error) {
	acl := jwt.MapClaims{}
	acl["authorized"] = true
	acl["exp"] = time.Now().Add(time.Hour * 6).Unix()
	acl["userId"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, acl)
	return token.SignedString(config.SecretKey) //secret
}
