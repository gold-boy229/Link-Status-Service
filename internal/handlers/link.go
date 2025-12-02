package handlers

import "github.com/labstack/echo"

type linkHandler struct{}

func NewLinkHandler() *linkHandler {
	return &linkHandler{}
}

func (h linkHandler) ProcessList(c echo.Context) error {
	return nil
}

func (h linkHandler) BuildPDF(c echo.Context) error {
	return nil
}
