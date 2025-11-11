package utils

import (
	"net/http"

	"github.com/astaxie/beego/context"
)

type HttpContext struct {
	*context.Context
}

type HttpRequest struct {
	*http.Request
}

func (httpContext *HttpContext) GetAuthClientId() int {
	tokenString := httpContext.GetCookie("auth")
	return parseAuthToken(tokenString).ClientId
}

func (httpRequest *HttpRequest) GetAuthClientId() int {
	token, _ := httpRequest.Cookie("auth")
	return parseAuthToken(token.Value).ClientId
}

func (httpContext *HttpContext) GetAuthLogin() string {
	tokenString := httpContext.GetCookie("auth")
	return parseAuthToken(tokenString).Login
}

func (httpRequest *HttpRequest) GetAuthLogin() string {
	token, _ := httpRequest.Cookie("auth")
	return parseAuthToken(token.Value).Login
}

func IsErrorWriteHeader(err error, w http.ResponseWriter, httpStatusCode int, message string) bool {
	if err == nil {
		return false
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
	return true
}

func parseAuthToken(token string) UserClaims {
	user, _ := CurrentUserJWT(token)
	return user
}

func (httpRequest *HttpRequest) GetAuthUserId() int {
	token, _ := httpRequest.Cookie("auth")
	return parseAuthToken(token.Value).UserId
}

func (httpContext *HttpContext) GetAuthUserId() int {
	tokenString := httpContext.GetCookie("auth")
	return parseAuthToken(tokenString).UserId
}
