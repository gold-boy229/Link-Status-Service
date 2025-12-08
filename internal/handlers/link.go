package handlers

type linkHandler struct {
	LinkService linkService
	PDFBuilder  pdfBuilder
}

func NewLinkHandler(linkService linkService, pdfBuilder pdfBuilder) *linkHandler {
	return &linkHandler{LinkService: linkService, PDFBuilder: pdfBuilder}
}
