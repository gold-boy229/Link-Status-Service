package handlers

import (
	"context"

	"Link-Status-Service/internal/entity"

	"codeberg.org/go-pdf/fpdf"
)

type linkService interface {
	GetStatus(context.Context, entity.LinkGetStatusParams) (entity.LinkGetStatusResult, error)
	GetStatusesOfLinkSets(context.Context, entity.LinkBuildPDSParams) (entity.LinkBuildPDSResult, error)
}

type pdfBuilder interface {
	BuildPDF(linkStatuses []entity.LinkStatus) *fpdf.Fpdf
}
