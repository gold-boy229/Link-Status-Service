package app

import "github.com/labstack/echo"

type linkHandlerI interface {
	ProcessList(c echo.Context) error
	BuildPDF(c echo.Context) error
}
