package postgres

import (
	"context"
	"github.com/jackc/pgx"
	"strconv"
	"thread-service/db"
	"thread-service/forum"
	"thread-service/models"
	"time"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository() *ThreadRepository {
	return &ThreadRepository{db.GetDB()}
}

type Thread struct {
	Author     string
	Slug       *string
	Votes      int
	Title      string
	Created    time.Time
	ForumName  string
	Id         int
	Message    string
	PostsCount int
}

func toPostgresThread(thread *models.Thread) *Thread {
	return &Thread{
		Id:         thread.Id,
		Author:     thread.Author,
		Slug:       thread.Slug,
		Votes:      thread.Votes,
		Title:      thread.Title,
		Created:    thread.Created,
		ForumName:  thread.ForumName,
		Message:    thread.Message,
		PostsCount: thread.PostsCount,
	}
}

func ToModelThread(thread *Thread) *models.Thread {
	return &models.Thread{
		Id:         thread.Id,
		Author:     thread.Author,
		Slug:       thread.Slug,
		Votes:      thread.Votes,
		Title:      thread.Title,
		Created:    thread.Created,
		ForumName:  thread.ForumName,
		Message:    thread.Message,
		PostsCount: thread.PostsCount,
	}
}

func (r *ThreadRepository) getThreadPostsCount(ctx context.Context, id int) (int, error) {
	postsCount := 0

	row := r.db.QueryRow("SELECT count(*) FROM posts WHERE tid = $1", id)
	err := row.Scan(&postsCount)

	if err != nil {
		return -1, err
	}

	return postsCount, nil
}

func (r *ThreadRepository) GetThreadBySlug(ctx context.Context, slug string) (*models.Thread, error) {
	row := r.db.QueryRow("SELECT * FROM threads WHERE slug = $1", slug)

	thread := Thread{}
	err := row.Scan(&thread.Author, &thread.Created, &thread.ForumName, &thread.Id,
		&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, forum.ErrThreadNotFound
	}

	postsCount, err := r.getThreadPostsCount(ctx, thread.Id)
	if err != nil {
		return nil, err
	}
	thread.PostsCount = postsCount
	// TODO temporary for tests
	thread.Created = thread.Created.Add(-3 * time.Hour)

	return ToModelThread(&thread), nil
}

func (r *ThreadRepository) GetThreadById(ctx context.Context, id int) (*models.Thread, error) {
	row := r.db.QueryRow("SELECT * FROM threads WHERE id = $1", id)

	thread := Thread{}
	err := row.Scan(&thread.Author, &thread.Created, &thread.ForumName, &thread.Id,
		&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, forum.ErrThreadNotFound
	}

	postsCount, err := r.getThreadPostsCount(ctx, thread.Id)
	if err != nil {
		return nil, err
	}
	thread.PostsCount = postsCount
	// TODO temporary for tests
	thread.Created = thread.Created.Add(-3 * time.Hour)

	return ToModelThread(&thread), nil
}

func (r *ThreadRepository) DeleteThread(ctx context.Context, id int) error {
	_, err := r.db.Exec("DELETE FROM threads WHERE id = $1", id)
	return err
}

func (r *ThreadRepository) ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error) {
	id, _ := strconv.Atoi(slug)

	var query string
	var row *pgx.Row

	if message == "" && title == "" {
		query = `SELECT * FROM threads WHERE id = $1 OR slug = $2`
		row = r.db.QueryRow(query, id, slug)
	} else if message == "" {
		query = `UPDATE threads  
				SET title = $1 
				WHERE id = $2 OR slug = $3 RETURNING *`
		row = r.db.QueryRow(query, title, id, slug)
	} else if title == "" {
		query = `UPDATE threads  
				SET message = $1
				WHERE id = $2 OR slug = $3 RETURNING *`
		row = r.db.QueryRow(query, message, id, slug)
	} else {
		query = `UPDATE threads  
				SET message = $1, title = $2 
				WHERE id = $3 OR slug = $4 RETURNING *`
		row = r.db.QueryRow(query, message, title, id, slug)
	}

	thread := Thread{}
	err := row.Scan(&thread.Author, &thread.Created, &thread.ForumName, &thread.Id,
		&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, forum.ErrThreadNotFound
	}
	// TODO temporary for tests
	thread.Created = thread.Created.Add(-3 * time.Hour)

	return ToModelThread(&thread), nil
}

func (r *ThreadRepository) checkThread(slug string) (int, error) {
	id, _ := strconv.Atoi(slug)

	err := r.db.QueryRow("SELECT id FROM threads WHERE slug = $1 OR id = $2", slug, id).Scan(&id)
	if err != nil {
		return -1, forum.ErrThreadNotFound
	}

	return id, nil
}

func (r *ThreadRepository) VoteForThread(ctx context.Context, slug string, vote *models.Vote) (*models.Thread, error) {
	id, _ := strconv.Atoi(slug)

	thread := Thread{}
	transaction, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	// check thread existence and get forum's slug
	err = transaction.QueryRow("SELECT author, created, id, forum, message, slug::text, title, votes "+
		"FROM threads WHERE slug = $1 OR id = $2", slug, id).Scan(
		&thread.Author, &thread.Created, &thread.Id, &thread.ForumName,
		&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		_ = transaction.Rollback()
		return nil, forum.ErrThreadNotFound
	}
	// create/update vote
	rows, err := transaction.Exec("UPDATE votes SET voice=$1 WHERE tid=$2 AND nickname=$3;",
		vote.Voice, thread.Id, vote.Nickname)
	if count := rows.RowsAffected(); count == 0 {
		_, err := transaction.Exec("INSERT INTO votes (nickname, tid, voice)"+
			"VALUES ($1, $2, $3);", vote.Nickname, thread.Id, vote.Voice)
		if err != nil {
			_ = transaction.Rollback()
			return nil, forum.ErrUserNotFound
		}
	}
	// get new votes
	err = transaction.QueryRow("SELECT votes FROM threads WHERE id = $1",
		thread.Id).Scan(&thread.Votes)
	if err != nil {
		_ = transaction.Rollback()
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}
	// TODO temporary for tests
	thread.Created = thread.Created.Add(-3 * time.Hour)

	return ToModelThread(&thread), err
}

func (r *ThreadRepository) CreateThread(ctx context.Context, slug string, t *models.Thread) (*models.Thread, error) {
	thread := toPostgresThread(t)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow("INSERT INTO threads (author, created, forum, message, title, slug) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", thread.Author, thread.Created,
		slug, thread.Message, thread.Title, thread.Slug).Scan(&thread.Id)
	if err != nil {
		_ = tx.Rollback()
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

		return ToModelThread(&conflictThread), forum.ErrThreadAlreadyExists
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return ToModelThread(thread), nil
}

func (r *ThreadRepository) formGettingThreadsQuery(slug, since string, limit int, sort bool) string {
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

func (r *ThreadRepository) GetForumThreads(ctx context.Context, slug, since string, limit, offset int, sort bool) ([]*models.Thread, error) {
	query := r.formGettingThreadsQuery(slug, since, limit, sort)

	var rows *pgx.Rows
	// TODO temporary for tests
	loc, _ := time.LoadLocation("Europe/Moscow")
	sinceTime, _ := time.ParseInLocation(time.RFC3339, since, loc)
	sinceTime = sinceTime.Add(3 * time.Hour)

	var err error
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
