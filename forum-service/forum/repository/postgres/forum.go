package postgres

import (
	"context"
	"forum-service/db"
	forumPkg "forum-service/forum"
	"forum-service/models"
	"github.com/jackc/pgx"
)

type Forum struct {
	Title   string
	Slug    string
	User    string
	Threads int
	Posts   int
}

type User struct {
	About    string
	Email    string
	Fullname string
	Nickname string
}

type ForumRepository struct {
	db *pgx.ConnPool
}

func NewForumRepository() *ForumRepository {
	return &ForumRepository{db.GetDB()}
}

func (r *ForumRepository) CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	f := toPostgresForum(forum)

	err := r.db.QueryRow("SELECT nickname FROM users WHERE nickname = $1", f.User).Scan(&f.User)
	if err != nil {
		return nil, forumPkg.ErrUserNotFound
	}

	_, err = r.db.Exec("INSERT INTO forums (slug, title, author) VALUES ($1, $2, $3)",
		f.Slug, f.Title, f.User)
	if err != nil {
		var conflictForum Forum
		row := r.db.QueryRow("SELECT slug, title, author, threads, posts FROM forums WHERE slug = $1",
			f.Slug)
		err = row.Scan(
			&conflictForum.Slug,
			&conflictForum.Title,
			&conflictForum.User,
			&conflictForum.Threads,
			&conflictForum.Posts,
		)

		if err != nil {
			return nil, err
		}
		return ToModelForum(&conflictForum), forumPkg.ErrForumAlreadyExists
	}

	return ToModelForum(f), nil
}

func (r *ForumRepository) GetForums(ctx context.Context) ([]*models.Forum, error) {
	rows, err := r.db.Query("SELECT posts, slug, threads, title, author FROM forums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forums []*models.Forum
	for rows.Next() {
		forum := Forum{}
		err = rows.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
		if err != nil {
			return nil, err
		}

		forums = append(forums, ToModelForum(&forum))
	}

	return forums, nil
}

func (r *ForumRepository) DeleteForum(ctx context.Context, slug string) error {
	_, err := r.db.Exec("DELETE FROM forums WHERE slug = $1", slug)
	return err
}

func (r *ForumRepository) GetForum(ctx context.Context, slug string) (*models.Forum, error) {
	row := r.db.QueryRow("SELECT posts, slug, threads, title, author FROM forums WHERE slug = $1", slug)

	var forum Forum
	err := row.Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	if err != nil {
		return nil, forumPkg.ErrForumNotFound
	}

	return ToModelForum(&forum), nil
}

func (r *ForumRepository) GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error) {
	var forum Forum
	err := r.db.QueryRow("SELECT slug FROM forums WHERE slug = $1", slug).Scan(&forum.Slug)
	if err != nil {
		return nil, forumPkg.ErrForumNotFound
	}
	// Form query. It' important to use fUser instead of nickname here, because
	// fUser's collation is overwrite to compare symbols like '-' or '.' correctly
	query := `SELECT about, email, fullname, fUser 
					FROM users JOIN forum_users ON fUser = nickname AND forum = $1 `

	if since != "" {
		if sort {
			query += "AND fUser < $2 "
		} else {
			query += "AND fUser > $2 "
		}
	}
	query += "ORDER BY fUser "
	if sort {
		query += "DESC "
	}
	if limit > 0 {
		if since != "" {
			query += "LIMIT $3"
		} else {
			query += "LIMIT $2"
		}
	}

	var rows *pgx.Rows
	if since != "" && limit > 0 {
		rows, err = r.db.Query(query, slug, since, limit)
	} else if since != "" {
		rows, err = r.db.Query(query, slug, since)
	} else if limit > 0 {
		rows, err = r.db.Query(query, slug, limit)
	} else {
		rows, err = r.db.Query(query, slug)
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var result []*models.User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		result = append(result, ToModelUser(&user))
	}
	return result, nil
}

func toPostgresForum(f *models.Forum) *Forum {
	return &Forum{
		Title:   f.Title,
		Slug:    f.Slug,
		User:    f.User,
		Threads: f.Threads,
		Posts:   f.Posts,
	}
}

func ToModelForum(f *Forum) *models.Forum {
	return &models.Forum{
		Title:   f.Title,
		Slug:    f.Slug,
		User:    f.User,
		Threads: f.Threads,
		Posts:   f.Posts,
	}
}

func ToModelUser(f *User) *models.User {
	return &models.User{
		About:    f.About,
		Nickname: f.Nickname,
		Fullname: f.Fullname,
		Email:    f.Email,
	}
}
