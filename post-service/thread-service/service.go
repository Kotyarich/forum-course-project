package thread_service

import (
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	"net/http"
)

type ThreadService struct {
	url string
	cb  *gobreaker.CircuitBreaker
}

func NewThreadService(url string) *ThreadService {
	var st gobreaker.Settings
	st.Name = "HTTP THREADS"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	return &ThreadService{
		url: url,
		cb:  gobreaker.NewCircuitBreaker(st),
	}
}

func (s *ThreadService) CheckForum(ctx context.Context, slug string) error {
	url := fmt.Sprintf("%s/thread/%s/details", s.url, slug)

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
