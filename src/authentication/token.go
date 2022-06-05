package authentication

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

//ValidateToken verify is received token is valid
func ValidateToken(r *http.Request) error {
	tokenString := getToken(r)
	token, err := jwt.Parse(tokenString, returnVerificationKey)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("Invalid Token!")
}

func getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func GetUserID(r *http.Request) (uint64, error) {
	tokenString := getToken(r)
	token, err := jwt.Parse(tokenString, returnVerificationKey)
	if err != nil {
		return 0, err
	}

	if permission, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", permission["userId"]), 10, 64) //base 10, 64 bits
		if err != nil {
			return 0, err
		}

		return userID, nil
	}

	return 0, errors.New("Invalid Token")
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected subscription method! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}
