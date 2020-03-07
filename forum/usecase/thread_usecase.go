package usecase

import (
	"context"
	"dbProject/forum"
	"dbProject/models"
	"strconv"
)

type ThreadUseCase struct {
	threadRepo  forum.RepositoryThread
}

func (u *ThreadUseCase) CreateThreadPost(ctx context.Context, slug string, posts []*models.Post) ([]*models.Post, error) {
	posts, err := u.threadRepo.ThreadPostCreate(ctx, slug, posts)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (u *ThreadUseCase) GetThread(ctx context.Context, slug string) (*models.Thread, error) {
	id, err := strconv.Atoi(slug)
	var thread *models.Thread
	if err != nil {
		thread, err = u.threadRepo.GetThreadBySlug(ctx, slug)
	} else {
		thread, err = u.threadRepo.GetThreadById(ctx, id)
	}

	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (u *ThreadUseCase) ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error) {
	thread, err := u.threadRepo.ChangeThread(ctx, slug, title, message)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (u *ThreadUseCase) GetThreadPosts(ctx context.Context, slug string, limit, since int, desc bool, sort models.PostSortType) ([]*models.Post, error) {
	var posts []*models.Post
	var err error

	switch sort {
	case models.Flat:
		posts, err = u.threadRepo.GetThreadPostsFlat(ctx, slug, limit, since, desc)
	case models.Tree:
		posts, err = u.threadRepo.GetThreadPostsTree(ctx, slug, limit, since, desc)
	case models.ParentTree:
		posts, err = u.threadRepo.GetThreadPostsParentTree(ctx, slug, limit, since, desc)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (u *ThreadUseCase) VoteForThread(ctx context.Context, slug string, vote models.Vote) (*models.Thread, error) {
	thread, err := u.threadRepo.VoteForThread(ctx, slug, &vote)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

