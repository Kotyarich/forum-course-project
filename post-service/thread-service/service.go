package thread_service

import (
	"context"
	"fmt"
	"net/http"
)

type ThreadService struct {
	url string
}

func NewThreadService(url string) *ThreadService {
	return &ThreadService{
		url: url,
	}
}

func (s *ThreadService) CheckForum(ctx context.Context, slug string) error {
	url := fmt.Sprintf("%s/thread/%s/details", s.url, slug)
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