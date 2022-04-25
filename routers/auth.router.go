package routers

import (
	"confessit/handlers"
	"confessit/models"
	"confessit/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate = validator.New()

func (r *Route) Signup(c echo.Context) error {
	var payload models.UserCreatePayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	hasUser := handlers.GetUser(r.Conn, payload.Username)

	if hasUser.Username != "" {
		return c.JSON(http.StatusBadRequest, "user already exist.")
	}

	if err := handlers.CreateUser(r.Conn, payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, "New user created.")
}

func (r *Route) Login(c echo.Context) error {
	var payload models.UserLoginPayload

	if err := (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := validate.Struct(payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	hasUser := handlers.GetUser(r.Conn, payload.Username)

	if hasUser.Username == "" {
		return c.JSON(http.StatusUnauthorized, "user does not exist.")
	}

	if err := hasUser.ValidatePassword(payload.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, "password mismatch.")
	}

	sessionToken := uuid.New()
	expiresAt := time.Now().Add(120 * time.Second)

	cookie := utils.CreateCookie("session_token", sessionToken.String(), expiresAt.Minute())

	var sessionCookie models.Session = models.Session{
		ID:       sessionToken,
		Username: payload.Username,
		Expiry:   expiresAt,
	}

	r.Conn.Create(&sessionCookie)

	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, "logged in")
}

func (r *Route) Refresh(c echo.Context) error {
	session, err := c.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return c.JSON(http.StatusUnauthorized, "no cookie in headers")
		}
		return c.JSON(http.StatusBadRequest, "err getting cookie in headers")
	}

	sessionToken := session.Value
	var token models.Session

	r.Conn.Model(&models.Session{}).Where("id = ?", sessionToken).Find(&token)

	if token.ID == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	if token.IsExpired() {
		new_cookie := utils.CreateCookie("session_token", "", -1)
		c.SetCookie(&new_cookie)
		r.Conn.Delete(&models.Session{}, "id = ?", sessionToken)

		return c.JSON(http.StatusUnauthorized, "expired cookie")
	}

	new_sessionToken := uuid.New()
	expiresAt := time.Now().Add(120 * time.Second)

	cookie := utils.CreateCookie("session_token", new_sessionToken.String(), expiresAt.Minute())

	var sessionCookie models.Session = models.Session{
		ID:       new_sessionToken,
		Username: token.Username,
		Expiry:   expiresAt,
	}

	r.Conn.Create(&sessionCookie)
	r.Conn.Delete(&models.Session{}, "id = ?", sessionToken)
	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, "refreshed")
}

func (r *Route) Logout(c echo.Context) error {
	session, err := c.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return c.JSON(http.StatusUnauthorized, "no cookie in headers")
		}
		return c.JSON(http.StatusBadRequest, "err getting cookie in headers")
	}

	sessionToken := session.Value
	r.Conn.Delete(&models.Session{}, "id = ?", sessionToken)
	nil_cookie := utils.CreateCookie("session_token", "", -1)

	c.SetCookie(&nil_cookie)
	return c.JSON(http.StatusOK, "logged out")
}
