package restful

import (
	"context"
	"io"
	"net/http"
	"time"
)

type IRestfulService interface {
	Get(ctx context.Context, url string, timeOut string) ([]byte, error)
}

type restfulService struct{}

func NewRestfulService() IRestfulService {
	return &restfulService{}
}

func (rs *restfulService) Get(ctx context.Context, url string, timeOut string) ([]byte, error) {
	to, _ := time.ParseDuration(timeOut)
	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
