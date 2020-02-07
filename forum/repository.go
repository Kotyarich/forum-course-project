package forum

import (
	"context"
	"dbProject/models"
)

type RepositoryForum interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	CreateThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error)
	GetForum(ctx context.Context, slug string) (*models.Forum, error)
	GetForumThreads(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.Thread, error)
	GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error)
}

type RepositoryThread interface {
	ThreadPostCreate(ctx context.Context, slug string, posts []*models.Post) ([]*models.Post, error)
	GetThreadBySlug(ctx context.Context, slug string) (*models.Thread, error)
	GetThreadById(ctx context.Context, id int) (*models.Thread, error)
	ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error)
	GetThreadPostsFlat(ctx context.Context, slug string, limit, since int, desc bool) ([]*models.Post, error)
	GetThreadPostsTree(ctx context.Context, slug string, limit, since int, desc bool) ([]*models.Post, error)
	GetThreadPostsParentTree(ctx context.Context, slug string, limit, since int, desc bool) ([]*models.Post, error)
	VoteForThread(ctx context.Context, slug string, vote *models.Vote) (*models.Thread, error)
}

type RepositoryPost interface {
	GetPostAuthor(ctx context.Context, nickname string) (*models.User, error)
	GetPostForum(ctx context.Context, slug string) (*models.Forum, error)
	GetPostThread(ctx context.Context, id int) (*models.Thread, error)
	GetPost(ctx context.Context, id int) (*models.Post, error)
	ChangePost(ctx context.Context, newMessage string, post *models.Post) error
}

type RepositoryService interface {
	Clear(ctx context.Context) error
	Status(ctx context.Context) (*models.Status, error)
}