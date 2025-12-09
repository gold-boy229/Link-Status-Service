package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const requestTimeout time.Duration = 2 * time.Second

type customHTTPClient struct {
	http.Client
}

func NewCustomHTTPClient() *customHTTPClient {
	return &customHTTPClient{Client: http.Client{Timeout: requestTimeout}}
}

func (client *customHTTPClient) IsLinkAvailable(ctx context.Context, link string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, link, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create HEAD request for link %q; err: %w", link, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("network error during HEAD request to link %q; err: %w", link, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("cannot close response body: %v", closeErr)
		}
	}()

	if isCodeSuccessful(resp.StatusCode) || isCodeRedirection(resp.StatusCode) {
		return true, nil
	}
	return false, nil
}

func isCodeSuccessful(code int) bool {
	return 200 <= code && code < 300
}

func isCodeRedirection(code int) bool {
	return 300 <= code && code < 400
}
