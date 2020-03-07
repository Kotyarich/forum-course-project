package usecase

import (
	"context"
	"dbProject/forum"
	"dbProject/models"
	"time"
)

type PostUseCase struct {
	postRepo forum.RepositoryPost
}

func (u *PostUseCase) GetPostInfo(ctx context.Context, id int, user, thread, forum bool) (*models.DetailedInfo, error) {
	details := new(models.DetailedInfo)

	post, err := u.postRepo.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	details.PostInfo = *post

	if user {
		details.AuthorInfo, err = u.postRepo.GetPostAuthor(ctx, post.Author)
		if err != nil {
			return nil, err
		}
	}
	if forum {
		details.ForumInfo, err = u.postRepo.GetPostForum(ctx, post.ForumName)
		if err != nil {
			return nil, err
		}
	}
	if thread {
		details.ThreadInfo, err = u.postRepo.GetPostThread(ctx, post.Tid)
		if err != nil {
			return nil, err
		}
		// TODO temporary for tests
		details.ThreadInfo.Created = details.ThreadInfo.Created.Add(-3 * time.Hour)
	}

	return details, nil
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
