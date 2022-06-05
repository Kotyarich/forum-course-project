package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"user-service/models"
	"user-service/auth"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Информация о пользователе
// swagger:model User
type userInput struct {
	// Описание пользователя
	//
	// example: This is the day you will always remember as the day that you almost caught Captain Jack Sparrow!
	About    string `json:"about"`

	// Почтовый адрес пользователя
	//
	// format: email
	// example: captaina@blackpearl.sea
	Email    string `json:"email"`

	// Полное имя пользователя
	//
	// example: Captain Jack Sparrow
	Fullname string `json:"fullname"`

	// Имя пользователя (уникальное поле)
	//
	// format: identity
	// read only: true
	// example: j.sparrow
	Nickname string `json:"nickname"`

	// Пароль пользователя
	//
	// example: 123456
	Password string `json:"password,omitempty"`
}

func userInputToModel(user userInput) *models.User {
	return &models.User{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Password: user.Password,
		Email:    user.Email,
		About:    user.About,
	}
}

// Данные пользователя для авторизации
// swagger:model UserLogIn
type signInInput struct {
	// Имя пользователя
	Nickname string `json:"nickname"`

	// Пароль пользователя
	Password string `json:"password"`
}

func (h *Handler) SignOutHandler(c echo.Context) error {
	_, err := c.Cookie("Auth")
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	c.SetCookie(&http.Cookie{
		Name:     "Auth",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	} else {
		return c.String(http.StatusOK, "")
	}
}

func (h *Handler) UserAuthHandler(c echo.Context) error {
	var input signInInput
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	ctx := context.WithValue(c.Request().Context(), "UserAgent", c.Request().UserAgent())

	u, token, err := h.useCase.SignIn(ctx, input.Nickname, input.Password)
	if err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}

	cookie := &http.Cookie{
		Name:     "Auth",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
		Domain:   "localhost",
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, UserToUserOutput(u))
}

func (h *Handler) UserCheckAuthHandler(c echo.Context) error {
	u := c.Request().Context().Value("auth")
	if u == nil {
		msg := map[string]string{"error": "not authorised"}
		return c.JSON(http.StatusForbidden, msg)
	}

	return c.JSON(http.StatusOK, UserToUserOutput(u.(*models.User)))
}

func (h *Handler) UserCreateHandler(c echo.Context) error {
	if c.Request().Method == http.MethodOptions {
		c.Request().Header.Set("content-type", "application/json")
		return c.String(http.StatusOK, "")
	}

	u := userInput{}
	if err := c.Bind(&u); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	model := userInputToModel(u)
	model.Nickname = c.Param("nickname")
	ctx := context.WithValue(c.Request().Context(), "UserAgent", c.Request().UserAgent())
	conflicts, token, err := h.useCase.SignUp(ctx, model)
	if err != nil {
		if conflicts == nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		c.Request().Header.Set("content-type", "application/json")
		return c.JSON(http.StatusConflict, UsersToJsonArray(conflicts))
	}

	cookie := &http.Cookie{
		Name:     "Auth",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
		Domain:   "localhost",
	}
	c.SetCookie(cookie)

	c.Request().Header.Set("content-type", "application/json")
	return c.JSON(http.StatusCreated, UserToUserOutput(model))
}

type UserOutput struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"isAdmin"`
}

func UserToUserOutput(user *models.User) *UserOutput {
	return &UserOutput{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Email:    user.Email,
		About:    user.About,
		IsAdmin:  user.IsAdmin,
	}
}

func UsersToJsonArray(users []*models.User) []*UserOutput {
	var result []*UserOutput
	for i := 0; i < len(users); i++ {
		result = append(result, UserToUserOutput(users[i]))
	}

	return result
}
