package forum

import "thread-service/models"

type AuthService interface {
	CheckAuth(token string) (*models.User, error)
}
