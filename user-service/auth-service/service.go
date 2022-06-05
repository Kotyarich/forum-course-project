package auth_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"user-service/models"
)

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

func usersInputToArray(users []userInput) []*models.User {
	result := make([]*models.User, 0, len(users))

	for _, u := range users {
		result = append(result, userInputToModel(u))
	}

	return result
}

func (s *AuthService) CheckAuth(ctx context.Context, token string) (*models.User, error) {
	url := fmt.Sprintf("%s/user/check", s.url)
	bearer := "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", bearer)

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
