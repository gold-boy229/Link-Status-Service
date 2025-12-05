package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockLinkChecker struct {
	mock.Mock
}

func NewMockLinkChecker() *mockLinkChecker {
	return &mockLinkChecker{}
}

func (m *mockLinkChecker) IsLinkAvailable(ctx context.Context, link string) (bool, error) {
	args := m.Called(ctx, link)
	return args.Bool(0), args.Error(1)
}
