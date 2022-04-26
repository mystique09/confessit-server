package utils

import (
	"net/http"
	"os"
)

func CreateCookie(name string, value string, maxAge int) http.Cookie {
	var mode string = os.Getenv("MODE")
	domain := os.Getenv("COOKIE_DOMAIN")

	new_cookie := new(http.Cookie)
	new_cookie.Name = name
	new_cookie.Value = value
	new_cookie.MaxAge = maxAge
	new_cookie.Path = "/"
	new_cookie.Domain = domain
	new_cookie.SameSite = http.SameSiteNoneMode
	new_cookie.HttpOnly = true
	new_cookie.Secure = mode == "production"

	return *new_cookie
}
