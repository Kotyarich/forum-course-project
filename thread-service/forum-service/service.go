package forum_service

import (
	"context"
	"fmt"
	"net/http"
)

type ForumService struct {
	url string
}

func NewForumService(url string) *ForumService {
	return &ForumService{
		url: url,
	}
}

func (s *ForumService) CheckForum(ctx context.Context, slug string) error {
	url := fmt.Sprintf("%sforum/%s/details", s.url, slug)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func() {_ = resp.Body.Close()}()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	return nil
}