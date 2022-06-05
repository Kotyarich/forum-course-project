package usecase

import (
	"forum-service/forum"
)

func NewForumUseCase(forumRepo forum.RepositoryForum) forum.UseCase {
	return forum.UseCase{
		ForumUseCase:   &ForumUseCase{forumRepo: forumRepo},
	}
}
