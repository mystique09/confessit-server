package utils

import (
	"confessit/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	claims := &JwtClaims{
		payload.Id,
		payload.Username,
		jwt.StandardClaims{
			Id:        payload.Id.String(),
			Issuer:    payload.Username,
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	raw_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return raw_token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GetPayloadUsername(user *jwt.Token) string {
	var claims jwt.MapClaims = user.Claims.(jwt.MapClaims)
	username := claims["username"]
	return username.(string)
}

type JwtClaims struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.StandardClaims
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewError(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}
