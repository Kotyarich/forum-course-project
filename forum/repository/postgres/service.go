package postgres

import (
	"context"
	"dbProject/db"
	"dbProject/models"
	"github.com/jackc/pgx"
)

type ServiceRepository struct {
	db *pgx.ConnPool
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{db.GetDB()}
}

func (r *ServiceRepository) Clear(ctx context.Context) error {
	_, err := r.db.Exec(`
			TRUNCATE TABLE forum_users, votes, posts, threads, sessions, forums, users 
			  RESTART IDENTITY`)

	return err
}

func (r *ServiceRepository) Status(ctx context.Context) (*models.Status, error) {
	forums := 0
	err := r.db.QueryRow(`SELECT COUNT(*) FROM forums`).Scan(&forums)
	if err != nil {
		return nil, err
	}

	threads := 0
	err = r.db.QueryRow(`SELECT COUNT(*) FROM threads`).Scan(&threads)
	if err != nil {
		return nil, err
	}

	posts := 0
	err = r.db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&posts)
	if err != nil {
		return nil, err
	}

	users := 0
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&users)
	if err != nil {
		return nil, err
	}

	return &models.Status{
		Forums:  forums,
		Threads: threads,
		Posts:   posts,
		Users:   users,
	}, nil
}
