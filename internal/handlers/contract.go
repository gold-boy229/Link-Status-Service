package handlers

import (
	"Link-Status-Service/internal/entity"
	"context"
)

type linkService interface {
	GetStatus(context.Context, entity.LinkGetStatus_Params) (entity.LinkGetStatus_Result, error)
}
