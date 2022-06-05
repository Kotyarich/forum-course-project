package usecase

import (
	"context"
	"fmt"
	"post-service/forum"
	"post-service/models"
	"post-service/thread-service"
	"strconv"
)

type PostUseCase struct {
	threadService thread_service.ThreadService
	postRepo      forum.RepositoryPost
	producer      forum.Producer
}

func (u *PostUseCase) CreateThreadPost(ctx context.Context, threadId int, posts []*models.Post) ([]*models.Post, error) {
	err := u.threadService.CheckForum(ctx, fmt.Sprintf("%d", threadId))
	if err != nil {
		return nil, forum.ErrThreadNotFound
	}

	posts, err = u.postRepo.ThreadPostCreate(ctx, threadId, posts)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		u.producer.Produce(strconv.Itoa(post.Id))
	}

	return posts, err
}

func (u *PostUseCase) GetThreadPosts(ctx context.Context, threadId int, limit, offset, since int, desc bool, sort models.PostSortType) ([]*models.Post, error) {
	var posts []*models.Post
	var err error

	switch sort {
	case models.Flat:
		posts, err = u.postRepo.GetThreadPostsFlat(ctx, threadId, limit, offset, since, desc)
	case models.Tree:
		posts, err = u.postRepo.GetThreadPostsTree(ctx, threadId, limit, offset, since, desc)
	case models.ParentTree:
		posts, err = u.postRepo.GetThreadPostsParentTree(ctx, threadId, limit, offset, since, desc)
	}

	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (u *PostUseCase) ChangePost(ctx context.Context, id int, message string) (*models.Post, error) {
	post, err := u.postRepo.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	if message != "" && message != post.Message {
		err = u.postRepo.ChangePost(ctx, message, post)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}
