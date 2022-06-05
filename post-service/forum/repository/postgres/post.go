package postgres

import (
	"context"
	"github.com/jackc/pgx"
	"post-service/db"
	"post-service/forum"
	"post-service/models"
	"time"
)

type PostRepository struct {
	db *pgx.ConnPool
}

func NewPostRepository() *PostRepository {
	return &PostRepository{db.GetDB()}
}

type Post struct {
	Author    string
	Created   time.Time
	ForumName string
	Id        int
	IsEdited  bool
	Message   string
	Parent    int
	Tid       int
}

func toModelPost(p *Post) *models.Post {
	return &models.Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Id:        p.Id,
		IsEdited:  p.IsEdited,
		Message:   p.Message,
		Parent:    p.Parent,
		Tid:       p.Tid,
	}
}

func (r *PostRepository) ThreadPostCreate(ctx context.Context, tid int, posts []*models.Post) ([]*models.Post, error) {
	if len(posts) != 0 {
		if posts[0].Parent != 0 {
			var parentTId int
			// Ignore error here because it will be handled with next if
			_ = r.db.QueryRow("SELECT tid FROM posts WHERE id = $1", posts[0].Parent).Scan(&parentTId)
			if parentTId != tid {
				return nil, forum.ErrWrongParentsThread
			}
		}
	}

	var err error
	for i := 0; i < len(posts); i++ {
		posts[i].Tid = tid
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

func (r *PostRepository) createPost(post *models.Post) error {
	query := "INSERT INTO posts (author, forum, message, parent, tid, slug, rootId) " +
		"VALUES ($1, $2, $3, $4, $5, " +
		"(SELECT slug FROM posts WHERE id = $4) || (SELECT currval('posts_id_seq')::integer), "
	if post.Parent == 0 {
		query += "(SELECT currval('posts_id_seq')::integer)) RETURNING id, created"
	} else {
		query += "(SELECT rootId FROM posts WHERE id = $4)) RETURNING id, created"
	}

	err := r.db.QueryRow(query,
		post.Author, post.ForumName, post.Message,
		post.Parent, post.Tid).Scan(&post.Id, &post.Created)

	if err != nil {
		return forum.ErrPostPatentNotFound
	}

	return nil
}

func (r *PostRepository) GetPost(ctx context.Context, id int) (*models.Post, error) {
	row := r.db.QueryRow(`SELECT author, created, forum, id, message, tid, isEdited, parent 
			FROM posts WHERE id = $1 `, id)

	post := Post{}
	err := row.Scan(&post.Author, &post.Created, &post.ForumName, &post.Id,
		&post.Message, &post.Tid, &post.IsEdited, &post.Parent)
	if err != nil {
		return nil, forum.ErrPostNotFound
	}

	return toModelPost(&post), nil
}

func (r *PostRepository) ChangePost(ctx context.Context, newMessage string, post *models.Post) error {
	_, err := r.db.Exec(`UPDATE posts SET message = $1, isEdited = TRUE WHERE id = $2`,
		newMessage, post.Id)
	if err != nil {
		return err
	}

	post.IsEdited = true
	post.Message = newMessage

	return nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
	return err
}


func (r *PostRepository) getPostsRows(query string, threadId, since, limit, offset int) (*pgx.Rows, error) {
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

func (r *PostRepository) postRowsToModelsArray(rows *pgx.Rows) ([]*models.Post, error) {
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

func (r *PostRepository) GetThreadPostsFlat(ctx context.Context, threadId int, limit, offset, since int, desc bool) ([]*models.Post, error) {
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

func (r *PostRepository) GetThreadPostsTree(ctx context.Context, threadId int, limit, offset, since int, desc bool) ([]*models.Post, error) {
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

func (r *PostRepository) GetThreadPostsParentTree(ctx context.Context, threadId int, limit, offset, since int, desc bool) ([]*models.Post, error) {
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
