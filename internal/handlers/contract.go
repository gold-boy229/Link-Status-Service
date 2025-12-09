package handlers

import (
	"Link-Status-Service/internal/entity"
	"context"

	"codeberg.org/go-pdf/fpdf"
)

type linkService interface {
	GetStatus(context.Context, entity.LinkGetStatusParams) (entity.LinkGetStatusResult, error)
	GetStatusesOfLinkSets(context.Context, entity.LinkBuildPDSParams) (entity.LinkBuildPDSResult, error)
}

type pdfBuilder interface {
	BuildPDF(linkStatuses []entity.LinkStatus) *fpdf.Fpdf
}
