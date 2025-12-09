package mocks

import (
	"Link-Status-Service/internal/entity"
	"context"

	"github.com/stretchr/testify/mock"
)

// must implement linkService interface
type mockLinkService struct {
	mock.Mock
}

func NewMockLinkService() *mockLinkService {
	return &mockLinkService{}
}

func (m *mockLinkService) GetStatus(ctx context.Context, params entity.LinkGetStatusParams) (entity.LinkGetStatusResult, error) {
	args := m.Called(ctx, params)

	result, _ := args.Get(0).(entity.LinkGetStatusResult)
	err := args.Error(1)
	return result, err
}

func (m *mockLinkService) GetStatusesOfLinkSets(ctx context.Context, params entity.LinkBuildPDSParams) (entity.LinkBuildPDSResult, error) {
	args := m.Called(ctx, params)

	result, _ := args.Get(0).(entity.LinkBuildPDSResult)
	err := args.Error(1)
	return result, err
}
