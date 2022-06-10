package forum_service

import (
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	"net/http"
)

type ForumService struct {
	url string
	cb  *gobreaker.CircuitBreaker
}

func NewForumService(url string) *ForumService {
	var st gobreaker.Settings
	st.Name = "HTTP AUTH"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	return &ForumService{
		url: url,
		cb:  gobreaker.NewCircuitBreaker(st),
	}
}

func (s *ForumService) CheckForum(ctx context.Context, slug string) error {
	url := fmt.Sprintf("%sforum/%s/details", s.url, slug)

	respI, err := s.cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		return resp, err
	})
	if err != nil {
		return err
	}

	resp := respI.(*http.Response)

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	return nil
}
