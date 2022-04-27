package utils

import (
	"confessit/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
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

func CreateJwt(payload models.JwtUserPayload) (string, error) {
	var claims JwtClaims = JwtClaims{
		payload,
		jwt.StandardClaims{
			Id:        payload.Id.String(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	raw_token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return raw_token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

type JwtClaims struct {
	Payload models.JwtUserPayload `json:"payload"`
	jwt.StandardClaims
}
