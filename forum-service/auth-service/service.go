package auth_service

import (
	"encoding/json"
	"fmt"
	"forum-service/models"
	"github.com/sony/gobreaker"
	"net/http"
	"time"
)

const CoockieName = "Auth"

type AuthService struct {
	url string
	cb  *gobreaker.CircuitBreaker
}

func NewAuthService(url string) *AuthService {
	var st gobreaker.Settings
	st.Name = "HTTP AUTH"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	return &AuthService{
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

	respI, err := s.cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("content-type", "application/json")
		req.AddCookie(&http.Cookie{Name: CoockieName, Value: token, Expires: time.Now().Add(time.Hour)})

		client := &http.Client{}
		resp, err := client.Do(req)
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
