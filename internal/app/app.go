package app

import (
	"Link-Status-Service/internal/handlers"

	"github.com/labstack/echo"
)

type app struct {
	echo *echo.Echo
}

func NewApp() *app {
	e := echo.New()
	e.Validator = NewCustomValidator()
	return &app{
		echo: echo.New(),
	}
}

func (a *app) Run() {
	var linkHandler linkHandlerI = handlers.NewLinkHandler()

	a.echo.POST("/links/get_status", linkHandler.GetStatus)
	a.echo.GET("/links/pdf", linkHandler.BuildPDF)

	a.echo.Logger.Fatal(a.echo.Start(":8080"))
}
