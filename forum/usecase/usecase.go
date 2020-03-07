package usecase

import (
	"dbProject/forum"
)

func NewForumUseCase(
	forumRepo forum.RepositoryForum,
	threadRepo forum.RepositoryThread,
	postRepo forum.RepositoryPost,
	serviceRepo forum.RepositoryService) forum.UseCase {
	return forum.UseCase{
		ForumUseCase:   &ForumUseCase{forumRepo: forumRepo},
		ThreadUseCase:  &ThreadUseCase{threadRepo: threadRepo},
		PostUseCase:    &PostUseCase{postRepo: postRepo},
		ServiceUseCase: &ServiceUseCase{serviceRepo: serviceRepo},
	}
}
