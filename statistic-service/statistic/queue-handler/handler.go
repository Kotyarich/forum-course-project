package queue_handler

import (
	"context"
	"fmt"
	"log"
	"statistic-service/statistic"
	"strconv"
)

type QueueHandler struct {
	statisticRepo statistic.Repository
}

func NewQueueHandler(userRepo statistic.Repository) *QueueHandler {
	return &QueueHandler{
		statisticRepo: userRepo,
	}
}

func (h *QueueHandler) HandlePostCreation(ctx context.Context, message string) error {
	id, err := strconv.Atoi(message)
	if err != nil {
		return fmt.Errorf("%v: converting message to id error", err)
	}

	err = h.statisticRepo.CreatePostRecord(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: creating record error", err)
	}

	return nil
}

func (h *QueueHandler) HandleUserCreation(ctx context.Context, message string) error {
	err := h.statisticRepo.CreateUserRecord(ctx, message)
	if err != nil {
		return fmt.Errorf("%v: creating record error", err)
	}

	return nil
}

func (h *QueueHandler) HandleVoteCreation(ctx context.Context, message string) error {
	log.Println("vote", message)

	err := h.statisticRepo.CreateVoteRecord(ctx, message)
	if err != nil {
		return fmt.Errorf("%v: creating record error", err)
	}

	log.Println("vote ok", message)

	return nil
}

func (h *QueueHandler) HandleThreadCreation(ctx context.Context, message string) error {
	log.Println("thread", message)
	id, err := strconv.Atoi(message)
	if err != nil {
		return fmt.Errorf("%v: converting message to id error", err)
	}

	log.Println("thread", id)
	err = h.statisticRepo.CreateThreadRecord(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: creating record error", err)
	}

	log.Println("thread ok", id)
	return nil
}

func (h *QueueHandler) HandleForumCreation(ctx context.Context, message string) error {
	id, err := strconv.Atoi(message)
	if err != nil {
		return fmt.Errorf("%v: converting message to id error", err)
	}

	err = h.statisticRepo.CreateForumRecord(ctx, id)
	if err != nil {
		return fmt.Errorf("%v: creating record error", err)
	}

	return nil
}
