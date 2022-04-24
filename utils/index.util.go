package utils

import (
	"net/http"
)

func CreateCookie(name string, value string, maxAge int) http.Cookie {
	new_cookie := new(http.Cookie)
	new_cookie.Name = name
	new_cookie.Value = value
	new_cookie.MaxAge = maxAge
	new_cookie.Path = "/"
	new_cookie.HttpOnly = true
	new_cookie.Secure = true

	return *new_cookie
}
