package usecase

import (
	"context"
	"dbProject/forum"
	"dbProject/models"
	"strconv"
)

type ForumUseCase struct {
	forumRepo  forum.RepositoryForum
	threadRepo forum.RepositoryThread
}

func NewForumUseCase(forumRepo forum.RepositoryForum, threadRepo forum.RepositoryThread) *ForumUseCase {
	return &ForumUseCase{
		forumRepo:  forumRepo,
		threadRepo: threadRepo,
	}
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

func (u *ForumUseCase) CreateThreadPost(ctx context.Context, slug string, posts []*models.Post) ([]*models.Post, error) {
	posts, err := u.threadRepo.ThreadPostCreate(ctx, slug, posts)
	if err != nil {
		return nil, err
	}

	return posts, err
}

func (u *ForumUseCase) GetThread(ctx context.Context, slug string) (*models.Thread, error) {
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

func (u *ForumUseCase) ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error) {
	thread, err := u.threadRepo.ChangeThread(ctx, slug, title, message)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (u *ForumUseCase) GetThreadPosts(ctx context.Context, slug string, limit, since int, desc bool, sort models.PostSortType) ([]*models.Post, error) {
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

func (u *ForumUseCase) VoteForThread(ctx context.Context, slug string, vote models.Vote) (*models.Thread, error) {
	thread, err := u.threadRepo.VoteForThread(ctx, slug, &vote)
	if err != nil {
		return nil, err
	}

	return thread, nil
}
