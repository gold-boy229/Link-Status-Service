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

	result, _ := args.Get(0).(entity.LinkGetStatus_Result)
	err := args.Error(1)
	return result, err
}

func (m *mockLinkService) GetStatusesOfLinkSets(ctx context.Context, params entity.LinkBuildPDS_Params) (entity.LinkBuildPDS_Result, error) {
	args := m.Called(ctx, params)

	result, _ := args.Get(0).(entity.LinkBuildPDS_Result)
	err := args.Error(1)
	return result, err
}
