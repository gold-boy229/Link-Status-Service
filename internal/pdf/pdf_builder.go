package pdf

import (
	"Link-Status-Service/internal/entity"

	"codeberg.org/go-pdf/fpdf"
)

type pdfBuilder struct{}

func NewPDFBuilder() *pdfBuilder {
	return &pdfBuilder{}
}

const (
	lineHeight float64 = 10

	headerWidth    float64 = 60
	headerFontSize float64 = 16

	linkColumnWidth   float64 = 130
	statusColumnWidth float64 = 60

	tableCellFontSize float64 = 12
)

type pdfCellParams struct {
	w         float64
	h         float64
	txtStr    string
	borderStr string
	ln        int
	alignStr  string
	fill      bool
	link      int
	linkStr   string
}

func newColumnParams(width float64) pdfCellParams {
	cellParams := newBasicColumnParams()
	cellParams.w = width
	return cellParams
}

func newBasicColumnParams() pdfCellParams {
	return pdfCellParams{
		w:         0,
		h:         lineHeight,
		txtStr:    "",
		borderStr: "1",
		ln:        0,
		alignStr:  "C",
		fill:      false,
		link:      0,
		linkStr:   "",
	}
}

func (builder *pdfBuilder) BuildPDF(linkStatuses []entity.LinkStatus) *fpdf.Fpdf {
	pdf := fpdf.New("Portrait", "mm", "A4", "")
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", headerFontSize)
	pdf.Cell(headerWidth, lineHeight, "Links availablity")
	pdf.Ln(lineHeight)

	linkColumnParams := newColumnParams(linkColumnWidth)
	statusColumnParams := newColumnParams(statusColumnWidth)

	// Set Head of Table
	pdf.SetFont("Arial", "B", tableCellFontSize)
	pdfAddCell(pdf, "Link", linkColumnParams)
	pdfAddCell(pdf, "Status", statusColumnParams)
	pdf.Ln(lineHeight)

	// Build table
	pdf.SetFont("Arial", "", tableCellFontSize)
	for _, linkStatus := range linkStatuses {
		pdfAddCell(pdf, linkStatus.Address, linkColumnParams)
		pdfAddCell(pdf, linkStatus.Status, statusColumnParams)
		pdf.Ln(lineHeight)
	}
	return pdf
}

func pdfAddCell(pdf *fpdf.Fpdf, txtStr string, params pdfCellParams) {
	pdf.CellFormat(
		params.w,
		params.h,
		txtStr,
		params.borderStr,
		params.ln,
		params.alignStr,
		params.fill,
		params.link,
		params.linkStr,
	)
}
