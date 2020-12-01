package http

import (
	"context"
	"dbProject/models"
	"dbProject/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Handler struct {
	useCase user.UseCase
}

func NewHandler(useCase user.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type userInput struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname"`
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

type signInInput struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func (h *Handler) SignOutHandler(c echo.Context) error {
	cookie, err := c.Cookie("Auth")
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = h.useCase.SignOut(c.Request().Context(), cookie.Value)
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

	return c.JSON(http.StatusOK,UserToUserOutput(u))
}

func (h *Handler) UserCheckAuthHandler(c echo.Context) error {
	u := c.Request().Context().Value("user")
	if u == nil {
		msg := map[string]string{"error": "not authorised"}
		return c.JSON(http.StatusForbidden, msg)
	}

	return c.JSON(http.StatusOK, UserToUserOutput(u.(*models.User)))
}

func (h *Handler) UserGetHandler(c echo.Context) error {
	nickname := c.Param("nickname")
	u, err := h.useCase.GetProfile(c.Request().Context(), nickname)
	if err != nil {
		if err == user.ErrUserNotFound {
			msg := map[string]string{"message": "404"}
			return c.JSON(http.StatusNotFound, msg)
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, UserToUserOutput(u))
}

func (h *Handler) UserPostHandler(c echo.Context) error {
	nickname := c.Param("nickname")

	var input userInput
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	input.Nickname = nickname

	newUser, err := h.useCase.ChangeProfile(c.Request().Context(), userInputToModel(input))

	if err == user.ErrUserAlreadyExists {
		msg := map[string]string{"message": "conflict"}
		return c.JSON(http.StatusConflict, msg)
	} else if err == user.ErrUserNotFound {
		msg := map[string]string{"message": "User not found"}
		return c.JSON(http.StatusNotFound, msg)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, UserToUserOutput(newUser))
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
}

func UserToUserOutput(user *models.User) *UserOutput {
	return &UserOutput{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Email:    user.Email,
		About:    user.About,
	}
}

func UsersToJsonArray(users []*models.User) []*UserOutput {
	var result []*UserOutput
	for i := 0; i < len(users); i++ {
		result = append(result, UserToUserOutput(users[i]))
	}

	return result
}
