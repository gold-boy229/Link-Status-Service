package handlers

import (
	"Link-Status-Service/internal/entity"
	"context"

	"codeberg.org/go-pdf/fpdf"
)

type linkService interface {
	GetStatus(context.Context, entity.LinkGetStatus_Params) (entity.LinkGetStatus_Result, error)
	GetStatusesOfLinkSets(context.Context, entity.LinkBuildPDS_Params) (entity.LinkBuildPDS_Result, error)
}

type pdfBuilder interface {
	BuildPDF(linkStatuses []entity.LinkStatus) *fpdf.Fpdf
}
