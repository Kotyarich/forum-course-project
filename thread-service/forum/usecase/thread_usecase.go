package usecase

import (
	"context"
	"fmt"
	"strconv"
	"thread-service/forum"
	"thread-service/forum-service"
	"thread-service/models"
)

type ThreadUseCase struct {
	forumService forum_service.ForumService
	threadRepo   forum.RepositoryThread
	producer     forum.Producer
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

func (u *ThreadUseCase) VoteForThread(ctx context.Context, slug string, vote models.Vote) (*models.Thread, error) {
	thread, err := u.threadRepo.VoteForThread(ctx, slug, &vote)
	if err != nil {
		return nil, err
	}

	u.producer.ProduceNewVote(fmt.Sprintf("%s: %s", slug, vote.Nickname))

	return thread, nil
}

func (u *ThreadUseCase) CreateForumThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error) {
	err := u.forumService.CheckForum(ctx, slug)
	if err != nil {
		return nil, forum.ErrForumNotFound
	}

	newThread, err := u.threadRepo.CreateThread(ctx, slug, thread)
	if err != nil {
		return newThread, err
	}

	u.producer.ProduceNewThread(strconv.Itoa(newThread.Id))

	return newThread, nil
}

func (u *ThreadUseCase) GetForumThreads(ctx context.Context, slug, since string, limit, offset int, sort bool) ([]*models.Thread, error) {
	threads, err := u.threadRepo.GetForumThreads(ctx, slug, since, limit, offset, sort)
	if err != nil {
		return nil, err
	}

	return threads, nil
}
