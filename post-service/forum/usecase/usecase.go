package usecase

import (
	"post-service/forum"
	"post-service/thread-service"
)

func NewForumUseCase(postRepo forum.RepositoryPost, service thread_service.ThreadService, producer forum.Producer) forum.UseCase {
	return forum.UseCase{
		PostUseCase:    &PostUseCase{threadService: service, postRepo: postRepo, producer: producer},
	}
}
