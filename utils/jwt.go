package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var privateKey = []byte(GetEnv("TOKEN_TTL", "your key"))

type UserClaims struct {
	Login    string
	UserName string
	ClientId int
	RoleId   int
	Rules    []string
	UserId   int
}

func JWTPrivateKey() []byte {
	return privateKey
}

func GenerateJWT(userClaims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":    userClaims.Login,
		"username": userClaims.UserName,
		"clientId": userClaims.ClientId,
		"roleId":   userClaims.RoleId,
		"rules":    userClaims.Rules,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"userId":   userClaims.UserId,
	})

	return token.SignedString(privateKey)
}

func ValidateJWT(tokenString string) error {
	token, err := getToken(tokenString)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func CurrentUserJWT(tokenString string) (userClaims UserClaims, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	err = ValidateJWT(tokenString)
	if err != nil {
		return UserClaims{}, err
	}

	token, _ := getToken(tokenString)
	claims, _ := token.Claims.(jwt.MapClaims)

	userClaims = UserClaims{
		Login:    claims["login"].(string),
		UserName: claims["username"].(string),
		ClientId: int(claims["clientId"].(float64)),
		RoleId:   int(claims["roleId"].(float64)),
		UserId:   int(claims["userId"].(float64)),
		Rules:    convertToStringSlice(claims["rules"]),
	}

	return userClaims, nil
}

func getToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func convertToStringSlice(rulesInterface interface{}) []string {
	rulesSlice, ok := rulesInterface.([]interface{})
	if !ok {
		return nil // or handle the error appropriately
	}

	rules := make([]string, len(rulesSlice))
	for i, rule := range rulesSlice {
		rules[i], ok = rule.(string)
		if !ok {
			return nil // or handle the error appropriately
		}
	}

	return rules
}
