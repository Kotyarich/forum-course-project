package auth_service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"post-service/models"
	"time"
)

const CoockieName = "Auth"

type AuthService struct {
	url string
}

func NewAuthService(url string) *AuthService {
	return &AuthService{
		url: url,
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

func (s *AuthService) CheckAuth(token string) (*models.User, error) {
	url := fmt.Sprintf("%suser/check", s.url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{Name: CoockieName, Value: token, Expires: time.Now().Add(time.Hour)})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {_ = resp.Body.Close()}()

	var user userInput
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return userInputToModel(user), nil
}