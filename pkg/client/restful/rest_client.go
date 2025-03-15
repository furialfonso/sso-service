package restful

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type IRestClient interface {
	Get(ctx context.Context, url string, timeOut string) ([]byte, error)
}

type restClient struct{}

func NewRestClient() IRestClient {
	return &restClient{}
}

func (rs *restClient) Get(ctx context.Context, url string, timeOut string) ([]byte, error) {
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

	if resp.StatusCode > 400 {
		code := strings.ReplaceAll(strings.ToLower(http.StatusText(resp.StatusCode)), " ", "_")
		return nil, fmt.Errorf("%d %s: %s", resp.StatusCode, code, string(body))
	}

	return body, nil
}
