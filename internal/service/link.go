package service

import (
	"Link-Status-Service/internal/entity"
	"context"
	"errors"
)

type linkService struct {
}

func NewLinkService() *linkService {
	return &linkService{}
}

func (s *linkService) GetStatus(ctx context.Context, params entity.LinkGetStatus_Params) (entity.LinkGetStatus_Result, error) {
	return entity.LinkGetStatus_Result{}, errors.New("not implemented")
}
