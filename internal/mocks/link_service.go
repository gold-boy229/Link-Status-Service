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

func (m *mockLinkService) GetStatus(ctx context.Context, params entity.LinkGetStatus_Params) (entity.LinkGetStatus_Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(entity.LinkGetStatus_Result), args.Error(1)
}

func (m *mockLinkService) GetStatusesOfLinkSets(ctx context.Context, params entity.LinkBuildPDS_Params) (entity.LinkBuildPDS_Result, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(entity.LinkBuildPDS_Result), args.Error(1)
}
