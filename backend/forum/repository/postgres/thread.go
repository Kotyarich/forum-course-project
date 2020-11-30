package postgres

import (
	"context"
	"dbProject/db"
	"dbProject/forum"
	"dbProject/models"
	"github.com/jackc/pgx"
	"strconv"
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

func (r *ThreadRepository) ThreadPostCreate(ctx context.Context, slug string, posts []*models.Post) ([]*models.Post, error) {
	tid, _ := strconv.Atoi(slug)
	post := Post{}

	err := r.db.QueryRow("SELECT id, forum "+
		"FROM threads WHERE slug = $1 OR id = $2", slug, tid).Scan(&post.Tid, &post.ForumName)
	if err != nil {
		return nil, forum.ErrThreadNotFound
	}

	if len(posts) != 0 {
		if posts[0].Parent != 0 {
			var parentTId int
			// Ignore error here because it will be handled with next if
			_ = r.db.QueryRow("SELECT tid FROM posts WHERE id = $1", posts[0].Parent).Scan(&parentTId)
			if parentTId != post.Tid {
				return nil, forum.ErrWrongParentsThread
			}
		}
	}

	for i := 0; i < len(posts); i++ {
		posts[i].Tid = post.Tid
		posts[i].ForumName = post.ForumName
		err = r.createPost(posts[i])
		if err != nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *ThreadRepository) createPost(post *models.Post) error {
	err := r.db.QueryRow("SELECT nickname "+
		"FROM users WHERE nickname = $1", post.Author).Scan(&post.Author)
	if err != nil {
		return forum.ErrUserNotFound
	}

	query := "INSERT INTO posts (author, forum, message, parent, tid, slug, rootId) " +
		"VALUES ($1, $2, $3, $4, $5, " +
		"(SELECT slug FROM posts WHERE id = $4) || (SELECT currval('posts_id_seq')::integer), "
	if post.Parent == 0 {
		query += "(SELECT currval('posts_id_seq')::integer)) RETURNING id, created"
	} else {
		query += "(SELECT rootId FROM posts WHERE id = $4)) RETURNING id, created"
	}

	err = r.db.QueryRow(query,
		post.Author, post.ForumName, post.Message,
		post.Parent, post.Tid).Scan(&post.Id, &post.Created)

	if err != nil {
		return forum.ErrPostPatentNotFound
	}

	return nil
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

func (r *ThreadRepository) getPostsRows(query string, threadId, since, limit, offset int) (*pgx.Rows, error) {
	var rows *pgx.Rows
	var err error
	if since >= 0 && limit > 0 {
		rows, err = r.db.Query(query, threadId, since, limit, offset)
	} else if since >= 0 {
		rows, err = r.db.Query(query, threadId, since, offset)
	} else if limit > 0 {
		rows, err = r.db.Query(query, threadId, limit, offset)
	} else {
		rows, err = r.db.Query(query, threadId, offset)
	}

	if err != nil {
		defer rows.Close()
		return nil, err
	}

	return rows, err
}

func (r *ThreadRepository) postRowsToModelsArray(rows *pgx.Rows) ([]*models.Post, error) {
	var result []*models.Post
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id, &post.Author, &post.Created, &post.ForumName,
			&post.IsEdited, &post.Message, &post.Parent, &post.Tid)
		if err != nil {
			return nil, err
		}

		result = append(result, toModelPost(&post))
	}

	return result, nil
}

func (r *ThreadRepository) checkThread(slug string) (int, error) {
	id, _ := strconv.Atoi(slug)

	err := r.db.QueryRow("SELECT id FROM threads WHERE slug = $1 OR id = $2", slug, id).Scan(&id)
	if err != nil {
		return -1, forum.ErrThreadNotFound
	}

	return id, nil
}

func (r *ThreadRepository) GetThreadPostsFlat(ctx context.Context, slug string, limit, offset, since int, desc bool) ([]*models.Post, error) {
	threadId, err := r.checkThread(slug)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, author, created, forum, isEdited, message, parent, tid " +
		"FROM posts WHERE tid = $1 "
	if since >= 0 {
		if desc {
			query += "AND id < $2 "
		} else {
			query += "AND id > $2 "
		}
	}
	query += "ORDER BY id "
	if desc {
		query += "DESC "
	}
	if limit > 0 {
		if since >= 0 {
			query += "LIMIT $3 OFFSET $4 "
		} else {
			query += "LIMIT $2 OFFSET $3"
		}
	} else {
		query += "OFFSET $2 "
	}

	rows, err := r.getPostsRows(query, threadId, since, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return r.postRowsToModelsArray(rows)
}

func (r *ThreadRepository) GetThreadPostsTree(ctx context.Context, slug string, limit, offset, since int, desc bool) ([]*models.Post, error) {
	threadId, err := r.checkThread(slug)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, author, created, forum, isEdited, message, parent, tid" +
		" FROM posts WHERE tid = $1 "
	if since >= 0 {
		if desc {
			query += "AND slug < (SELECT slug FROM posts WHERE id = $2) "
		} else {
			query += "AND slug > (SELECT slug FROM posts WHERE id = $2) "
		}
	}
	query += "ORDER BY slug "
	if desc {
		query += "DESC "
	}
	if limit > 0 {
		if since >= 0 {
			query += "LIMIT $3 OFFSET $4 "
		} else {
			query += "LIMIT $2 OFFSET $3 "
		}
	} else {
		query += "OFFSET $2 "
	}

	rows, err := r.getPostsRows(query, threadId, since, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.postRowsToModelsArray(rows)
}

func (r *ThreadRepository) GetThreadPostsParentTree(ctx context.Context, slug string, limit, offset, since int, desc bool) ([]*models.Post, error) {
	threadId, err := r.checkThread(slug)
	if err != nil {
		return nil, err
	}

	query := "WITH roots AS ( " +
		"SELECT id FROM posts WHERE tid = $1 AND parent = 0 "
	if since >= 0 {
		if desc {
			query += "AND id < (SELECT rootId FROM posts WHERE id = $2) "
		} else {
			query += "AND id > (SELECT rootId FROM posts WHERE id = $2) "
		}
	}
	query += "ORDER BY id "
	if desc {
		query += "DESC "
	}
	if limit > 0 {
		if since >= 0 {
			query += "LIMIT $3 OFFSET $4 "
		} else {
			query += "LIMIT $2 OFFSET $3 "
		}
	} else {
		query += "OFFSET $2"
	}
	query += ") SELECT posts.id, author, created, forum, isEdited, message, parent, tid " +
		"FROM posts JOIN roots ON roots.id = rootId "

	query += "ORDER BY "
	if desc {
		query += " rootId DESC, slug"
	} else {
		query += " slug"
	}

	rows, err := r.getPostsRows(query, threadId, since, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.postRowsToModelsArray(rows)
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
