package usecase

import (
	"context"
	"dbProject/forum"
	"dbProject/models"
)

type ForumUseCase struct {
	forumRepo forum.Repository
}

func NewForumUseCase(forumRepo forum.Repository) *ForumUseCase {
	return &ForumUseCase{forumRepo: forumRepo}
}

func (u *ForumUseCase) CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	newForum, err := u.forumRepo.CreateForum(ctx, forum)
	if err != nil {
		return newForum, err
	}

	return newForum, nil
}

func (u *ForumUseCase) CreateForumThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error) {
	newThread, err := u.forumRepo.CreateThread(ctx, slug, thread)
	if err != nil {
		return newThread, err
	}

	return newThread, nil
}

func (u *ForumUseCase) GetForumDetails(ctx context.Context, slug string) (*models.Forum, error) {
	f, err := u.forumRepo.GetForum(ctx, slug)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (u *ForumUseCase) GetForumThreads(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.Thread, error) {
	threads, err := u.forumRepo.GetForumThreads(ctx, slug, since, limit, sort)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (u *ForumUseCase) GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error) {
	users, err := u.forumRepo.GetForumUsers(ctx, slug, since, limit, sort)
	if err != nil {
		return nil, err
	}

	return users, nil
}
