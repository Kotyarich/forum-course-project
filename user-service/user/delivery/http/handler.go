package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"user-service/models"
	"user-service/user"
)

type Handler struct {
	useCase user.UseCase
}

func NewHandler(useCase user.UseCase) *Handler {
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
	About string `json:"about"`

	// Почтовый адрес пользователя
	//
	// format: email
	// example: captaina@blackpearl.sea
	Email string `json:"email"`

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

	IsAdmin bool `json:"isAdmin"`
}

func userInputToModel(user userInput) *models.User {
	return &models.User{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Password: user.Password,
		Email:    user.Email,
		About:    user.About,
		IsAdmin:  user.IsAdmin,
	}
}

func (s *Handler) checkAuth(token string) (*models.User, error) {
	url := fmt.Sprintf("%suser/check", "http://localhost:5002/")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{Name: "Auth", Value: token, Expires: time.Now().Add(time.Hour)})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	var user userInput
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return userInputToModel(user), nil
}

func (h *Handler) UserCheckHandler(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	u, err := h.useCase.CheckAuth(c.Request().Context(), username, password)
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

	cookie, err := c.Cookie("Auth")
	if err != nil {
		return c.String(http.StatusForbidden, "forbidden")
	}

	u, err := h.checkAuth(cookie.Value)
	if err != nil && u == nil || (u.Nickname != nickname && !u.IsAdmin) {
		return c.String(http.StatusForbidden, "forbidden")
	}

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

	//authed := c.Request().Context().Value("auth").(*models.User)
	//if authed == nil || u.Nickname != authed.Nickname {
	//	return c.String(http.StatusForbidden, "forbidden")
	//}

	ctx := context.WithValue(c.Request().Context(), "UserAgent", c.Request().UserAgent())
	conflicts, err := h.useCase.CreateUser(ctx, model)
	if err != nil {
		if conflicts == nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		fmt.Println(err.Error())
		c.Request().Header.Set("content-type", "application/json")
		return c.JSON(http.StatusConflict, UsersToJsonArray(conflicts))
	}

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
