package postgres

import (
	"context"
	"dbProject/db"
	"dbProject/forum"
	"dbProject/models"
	user "dbProject/user/repository/postgres"
	"github.com/jackc/pgx"
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

func modelToPost(p *models.Post) *Post {
	return &Post{
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

func (r *PostRepository) GetPostAuthor(ctx context.Context, nickname string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT about, email, fullname, nickname FROM users WHERE nickname = $1 `,
		nickname)

	author := user.User{}
	err := row.Scan(&author.About, &author.Email, &author.Fullname, &author.Nickname)
	if err != nil {
		return nil, forum.ErrUserNotFound
	}

	return user.ToModel(&author), nil
}

func (r *PostRepository) GetPostForum(ctx context.Context, slug string) (*models.Forum, error) {
	row := r.db.QueryRow(`SELECT posts, slug, threads, title, author FROM forums WHERE slug = $1 `,
		slug)

	postForum := Forum{}
	err := row.Scan(
		&postForum.Posts,
		&postForum.Slug,
		&postForum.Threads,
		&postForum.Title,
		&postForum.User,
	)
	if err != nil {
		return nil, forum.ErrForumNotFound
	}

	return ToModelForum(&postForum), nil
}

func (r *PostRepository) GetPostThread(ctx context.Context, id int) (*models.Thread, error) {
	row := r.db.QueryRow(`SELECT author, created, forum, id, message, slug, title, votes 
			FROM threads WHERE id = $1 `, id)

	thread := Thread{}
	err := row.Scan(&thread.Author, &thread.Created, &thread.ForumName, &thread.Id,
		&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, forum.ErrThreadNotFound
	}

	return ToModelThread(&thread), nil
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
