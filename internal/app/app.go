package app

import (
	"Link-Status-Service/internal/handlers"
	"Link-Status-Service/internal/repository"
	"Link-Status-Service/internal/service"
	"Link-Status-Service/internal/utils"

	"github.com/labstack/echo"
)

type app struct {
	echo *echo.Echo
}

func NewApp() *app {
	e := echo.New()
	e.Validator = utils.NewCustomValidator()
	return &app{
		echo: e,
	}
}

func (a *app) Run() {
	linkRepo := repository.NewLinkRepository()
	linkService := service.NewLinkService(linkRepo, &service.HTTPLinkChecker{})
	var linkHandler linkHandlerI = handlers.NewLinkHandler(linkService)

	a.echo.POST("/links/get_status", linkHandler.GetStatus)
	a.echo.GET("/links/pdf", linkHandler.BuildPDF)

	a.echo.Logger.Fatal(a.echo.Start(":8080"))
}
