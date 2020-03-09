package forum

import (
	"context"
	"dbProject/models"
)

type UseCaseForum interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	CreateForumThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error)
	GetForums(ctx context.Context) ([]*models.Forum, error)
	GetForumDetails(ctx context.Context, slug string) (*models.Forum, error)
	GetForumThreads(ctx context.Context, slug, since string, limit, offset int, sort bool) ([]*models.Thread, error)
	GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error)
	// TODO add removing
	// TODO add changing
}

type UseCaseThread interface {
	CreateThreadPost(ctx context.Context, slug string, posts []*models.Post) ([]*models.Post, error)
	GetThread(ctx context.Context, slug string) (*models.Thread, error)
	ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error)
	GetThreadPosts(ctx context.Context, slug string, limit, since int, desc bool, sort models.PostSortType) ([]*models.Post, error)
	VoteForThread(ctx context.Context, slug string, vote models.Vote) (*models.Thread, error)
}

type UseCasePost interface {
	GetPostInfo(ctx context.Context, id int, user, thread, forum bool) (*models.DetailedInfo, error) // extra
	ChangePost(ctx context.Context, id int, message string) (*models.Post, error)
	// TODO add removing
}

type UseCaseService interface {
	Clear(ctx context.Context) error
	Status(ctx context.Context) (*models.Status, error)
}

type UseCase struct {
	ForumUseCase   UseCaseForum
	ThreadUseCase  UseCaseThread
	PostUseCase    UseCasePost
	ServiceUseCase UseCaseService
}
