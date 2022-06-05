package forum

import "forum-service/models"

type AuthService interface {
	CheckAuth(token string) (*models.User, error)
}
