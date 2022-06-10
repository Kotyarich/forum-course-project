package user_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"net/http"
	"user-service/models"
)

type UserService struct {
	url string
	cb  *gobreaker.CircuitBreaker
}

func NewUserService(url string) *UserService {
	var st gobreaker.Settings
	st.Name = "HTTP AUTH"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	return &UserService{
		url: url,
		cb:  gobreaker.NewCircuitBreaker(st),
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

func usersInputToArray(users []userInput) []*models.User {
	result := make([]*models.User, 0, len(users))

	for _, u := range users {
		result = append(result, userInputToModel(u))
	}

	return result
}

func (s *UserService) CheckAuth(ctx context.Context, username string, password string) (*models.User, error) {
	url := fmt.Sprintf("%suser/check?username=%s&password=%s", s.url, username, password)
	respI, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		return resp, err
	})
	resp := respI.(*http.Response)
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

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, []*models.User, error) {
	url := fmt.Sprintf("%suser/%s/create", s.url, user.Nickname)

	jsonString, _ := json.Marshal(user)
	respI, err := s.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
		req.Header.Add("content-type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		return resp, err
	})
	if err != nil {
		return nil, nil, err
	}
	resp := respI.(*http.Response)

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusConflict {
		var users []userInput
		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			return nil, nil, err
		}
		return nil, usersInputToArray(users), nil
	}

	var created userInput
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, nil, err
	}
	return userInputToModel(created), nil, nil
}

func (s *UserService) GetUser(ctx context.Context, username string) (*models.User, error) {
	url := fmt.Sprintf("%suser/%s/profile", s.url, username)
	respI, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		return resp, err
	})
	if err != nil {
		return nil, err
	}
	resp := respI.(*http.Response)

	defer func() { _ = resp.Body.Close() }()

	var user userInput
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return userInputToModel(user), nil
}
