package pdf

import (
	"Link-Status-Service/internal/entity"

	"codeberg.org/go-pdf/fpdf"
)

type pdfBuilder struct{}

func NewPDFBuilder() *pdfBuilder {
	return &pdfBuilder{}
}

const LINE_HEIGHT float64 = 20

type pdf_cellParams struct {
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

var basicColumnParams = pdf_cellParams{
	w:         0,
	h:         10,
	txtStr:    "",
	borderStr: "1",
	ln:        0,
	alignStr:  "C",
	fill:      false,
	link:      0,
	linkStr:   "",
}

func (builder *pdfBuilder) BuildPDF(linkStatuses []entity.LinkStatus) *fpdf.Fpdf {
	pdf := fpdf.New("Portrait", "mm", "A4", "")
	pdf.AddPage()

	// Set Header
	pdf.SetFont("Arial", "B", 16) // is this line necessary?
	pdf.Cell(60, 10, "Links availablity")
	pdf.Ln(LINE_HEIGHT)

	linkColumnParams := basicColumnParams
	linkColumnParams.w = 90
	statusColumnParams := basicColumnParams
	statusColumnParams.w = 60

	// Set Head of Table
	pdf_addCell(pdf, "Link", linkColumnParams)
	pdf_addCell(pdf, "Status", statusColumnParams)
	pdf.Ln(LINE_HEIGHT)

	// Build table
	for _, linkStatus := range linkStatuses {
		pdf_addCell(pdf, linkStatus.Address, linkColumnParams)
		pdf_addCell(pdf, linkStatus.Status, statusColumnParams)
		pdf.Ln(LINE_HEIGHT)
	}
	return pdf
}

func pdf_addCell(pdf *fpdf.Fpdf, txtStr string, params pdf_cellParams) {
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
