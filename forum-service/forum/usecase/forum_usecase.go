package usecase

import (
	"context"
	"forum-service/forum"
	"forum-service/models"
)

type ForumUseCase struct {
	forumRepo forum.RepositoryForum
}

func (u *ForumUseCase) CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	newForum, err := u.forumRepo.CreateForum(ctx, forum)
	if err != nil {
		return newForum, err
	}

	return newForum, nil
}

func (u *ForumUseCase) GetForums(ctx context.Context) ([]*models.Forum, error) {
	forums, err := u.forumRepo.GetForums(ctx)
	if err != nil {
		return nil, err
	}

	return forums, nil
}

func (u *ForumUseCase) GetForumDetails(ctx context.Context, slug string) (*models.Forum, error) {
	f, err := u.forumRepo.GetForum(ctx, slug)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (u *ForumUseCase) GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error) {
	users, err := u.forumRepo.GetForumUsers(ctx, slug, since, limit, sort)
	if err != nil {
		return nil, err
	}

	return users, nil
}
