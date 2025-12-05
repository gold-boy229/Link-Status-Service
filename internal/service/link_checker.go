package service

import (
	"context"
	"errors"
)

type HTTPLinkChecker struct{}

func (h *HTTPLinkChecker) IsLinkAvailable(ctx context.Context, link string) (bool, error) {
	return false, errors.New("not implemented")
}
