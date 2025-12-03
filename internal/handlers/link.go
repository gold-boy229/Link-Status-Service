package handlers

type linkHandler struct {
	LinkService linkService
}

func NewLinkHandler(linkService linkService) *linkHandler {
	return &linkHandler{LinkService: linkService}
}
