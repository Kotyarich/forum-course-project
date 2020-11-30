package postgres

import (
	"context"
	"dbProject/db"
	forumPkg "dbProject/forum"
	"dbProject/models"
	"dbProject/user/repository/postgres"
	"github.com/jackc/pgx"
	"time"
)

type Forum struct {
	Title   string
	Slug    string
	User    string
	Threads int
	Posts   int
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

func (r *ForumRepository) CreateThread(ctx context.Context, slug string, t *models.Thread) (*models.Thread, error) {
	thread := toPostgresThread(t)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT nickname FROM users WHERE nickname = $1", thread.Author)
	err = row.Scan(&thread.Author)
	if err != nil {
		return nil, forumPkg.ErrUserNotFound
	}

	err = tx.QueryRow("SELECT slug FROM forums WHERE slug = $1", slug).Scan(&thread.ForumName)
	if err != nil {
		return nil, forumPkg.ErrForumNotFound
	}

	err = tx.QueryRow("INSERT INTO threads (author, created, forum, message, title, slug) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", thread.Author, thread.Created,
		thread.ForumName, thread.Message, thread.Title, thread.Slug).Scan(&thread.Id)
	if err != nil {
		tx.Rollback()
		row := r.db.QueryRow("SELECT * FROM threads WHERE slug = $1", thread.Slug)

		var conflictThread Thread
		err = row.Scan(
			&conflictThread.Author, &conflictThread.Created, &conflictThread.ForumName, &conflictThread.Id,
			&conflictThread.Message, &conflictThread.Slug, &conflictThread.Title, &conflictThread.Votes,
		)
		if err != nil {
			return nil, err
		}
		conflictThread.Created = conflictThread.Created.Add(-time.Hour * 3)

		return ToModelThread(&conflictThread), forumPkg.ErrThreadAlreadyExists
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return ToModelThread(thread), nil
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

func (r *ForumRepository) formGettingThreadsQuery(slug, since string, limit int, sort bool) string {
	query := "SELECT * FROM threads WHERE forum = $1 "
	if since != "" {
		if sort {
			query += "AND created <= $2 "
		} else {
			query += "AND created >= $2 "
		}
	}
	query += "ORDER BY created "
	if sort {
		query += "DESC "
	}
	if limit > 0 {
		if since != "" {
			query += "LIMIT $3 OFFSET $4"
		} else {
			query += "LIMIT $2 OFFSET $3"
		}
	} else {
		query += "OFFSET $2"
	}

	return query
}

func (r *ForumRepository) GetForumThreads(ctx context.Context, slug, since string, limit, offset int, sort bool) ([]*models.Thread, error) {
	var forum Forum
	err := r.db.QueryRow("SELECT slug FROM forums WHERE slug = $1", slug).Scan(&forum.Slug)
	if err != nil {
		return nil, forumPkg.ErrForumNotFound
	}

	query := r.formGettingThreadsQuery(slug, since, limit, sort)

	var rows *pgx.Rows
	// TODO temporary for tests
	loc, _ := time.LoadLocation("Europe/Moscow")
	sinceTime, _ := time.ParseInLocation(time.RFC3339, since, loc)
	sinceTime = sinceTime.Add(3 * time.Hour)

	if since != "" && limit > 0 {
		rows, err = r.db.Query(query, slug, sinceTime, limit, offset)
	} else if since != "" {
		rows, err = r.db.Query(query, slug, sinceTime, offset)
	} else if limit > 0 {
		rows, err = r.db.Query(query, slug, limit, offset)
	} else {
		rows, err = r.db.Query(query, slug, offset)
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var result []*models.Thread
	for rows.Next() {
		thr := Thread{}
		err = rows.Scan(&thr.Author, &thr.Created, &thr.ForumName, &thr.Id,
			&thr.Message, &thr.Slug, &thr.Title, &thr.Votes)
		if err != nil {
			return nil, err
		}
		result = append(result, ToModelThread(&thr))
	}
	return result, nil
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
		user := postgres.User{}
		err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			return nil, err
		}
		result = append(result, postgres.ToModel(&user))
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
