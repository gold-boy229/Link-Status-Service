package pdf

import (
	"Link-Status-Service/internal/entity"

	"codeberg.org/go-pdf/fpdf"
)

type pdfBuilder struct{}

func NewPDFBuilder() *pdfBuilder {
	return &pdfBuilder{}
}

const LineHeight float64 = 10

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

var basicColumnParams = pdfCellParams{
	w:         0,
	h:         LineHeight,
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
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(60, LineHeight, "Links availablity")
	pdf.Ln(LineHeight)

	linkColumnParams := basicColumnParams
	linkColumnParams.w = 130
	statusColumnParams := basicColumnParams
	statusColumnParams.w = 60

	// Set Head of Table
	pdf.SetFont("Arial", "B", 12)
	pdfAddCell(pdf, "Link", linkColumnParams)
	pdfAddCell(pdf, "Status", statusColumnParams)
	pdf.Ln(LineHeight)

	// Build table
	pdf.SetFont("Arial", "", 12)
	for _, linkStatus := range linkStatuses {
		pdfAddCell(pdf, linkStatus.Address, linkColumnParams)
		pdfAddCell(pdf, linkStatus.Status, statusColumnParams)
		pdf.Ln(LineHeight)
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
