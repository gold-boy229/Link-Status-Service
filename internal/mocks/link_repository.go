package mocks

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"
)

type mockLinkRepository struct {
	mock.Mock
}

func NewMockLinkRepository() *mockLinkRepository {
	return &mockLinkRepository{}
}

func (m *mockLinkRepository) GetLinkNum(ctx context.Context, links []string) (linkNum int, isNew bool, err error) {
	return 0, false, errors.New("not implemented")
}

func (m *mockLinkRepository) GetLinksByLinkNum(ctx context.Context, linkNum int) (links []string, err error) {
	args := m.Called(ctx, linkNum)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockLinkRepository) StoreLinks(ctx context.Context, links []string, linkNum int) error {
	return errors.New("not implemented")
}
