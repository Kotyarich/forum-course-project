package forum

import "post-service/models"

type AuthService interface {
	CheckAuth(token string) (*models.User, error)
}
