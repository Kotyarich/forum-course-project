package usecase

import (
	"thread-service/forum"
	"thread-service/forum-service"
)

func NewForumUseCase(threadRepo forum.RepositoryThread, service forum_service.ForumService, producer forum.Producer) forum.UseCase {
	return forum.UseCase{
		ThreadUseCase:  &ThreadUseCase{forumService: service, threadRepo: threadRepo, producer: producer},
	}
}
