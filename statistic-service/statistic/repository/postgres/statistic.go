package postgres

import (
	"context"
	"github.com/jackc/pgx"
	"statistic-service/db"
	"statistic-service/models"
)

type Statistic struct {
	Users   int64
	Posts   int64
	Votes   int64
	Threads int64
	Forums  int64
}

type UserRepository struct {
	db *pgx.ConnPool
}

func NewStatisticRepository() *UserRepository {
	return &UserRepository{db.GetDB()}
}

func (r UserRepository) CreateUserRecord(ctx context.Context, id string) error {
	_, err := r.db.Query("INSERT INTO users_creations (id, created_at) VALUES ($1, NOW())", id)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) CreatePostRecord(ctx context.Context, id int) error {
	_, err := r.db.Query("INSERT INTO posts_creations (id, created_at) VALUES ($1, NOW())", id)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) CreateVoteRecord(ctx context.Context, id string) error {
	_, err := r.db.Query("INSERT INTO votes_creations (id, created_at) VALUES ($1, NOW())", id)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) CreateThreadRecord(ctx context.Context, id int) error {
	_, err := r.db.Query("INSERT INTO threads_creations (id, created_at) VALUES ($1, NOW())", id)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) CreateForumRecord(ctx context.Context, id int) error {
	_, err := r.db.Query("INSERT INTO forums_creations (id, created_at) VALUES ($1, NOW())", id)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetStatistic(ctx context.Context) (*models.Status, error) {
	var statistic Statistic

	err := r.db.QueryRow(`
	SELECT COUNT(*)
		FROM posts_creations
	WHERE
		NOW() - created_at < '24 hour'::interval
	`).Scan(&statistic.Posts)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`
	SELECT COUNT(*)
		FROM users_creations
	WHERE
		NOW() - created_at < '24 hour'::interval
	`).Scan(&statistic.Users)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`
	SELECT COUNT(*)
		FROM threads_creations
	WHERE
		NOW() - created_at < '24 hour'::interval
	`).Scan(&statistic.Threads)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`
	SELECT COUNT(*)
		FROM votes_creations
	WHERE
		NOW() - created_at < '24 hour'::interval
	`).Scan(&statistic.Votes)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(`
	SELECT COUNT(*)
		FROM forums_creations
	WHERE
		NOW() - created_at < '24 hour'::interval
	`).Scan(&statistic.Forums)
	if err != nil {
		return nil, err
	}

	return ToModel(&statistic), nil
}

func ToModel(s *Statistic) *models.Status {
	return &models.Status{
		Users:   s.Users,
		Posts:   s.Posts,
		Votes:   s.Votes,
		Threads: s.Threads,
		Forums:  s.Forums,
	}
}
