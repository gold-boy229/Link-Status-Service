package app

import (
	"Link-Status-Service/internal/client"
	"Link-Status-Service/internal/handlers"
	"Link-Status-Service/internal/pdf"
	"Link-Status-Service/internal/repository"
	"Link-Status-Service/internal/service"
	"Link-Status-Service/internal/utils"
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	var dataManager dataManager = linkRepo
	if err := dataManager.LoadDataFromJSON(); err != nil {
		log.Fatalf("error during loading data from JSON file(s): %v\n", err)
	}

	clientChecker := client.NewCustomHTTPClient()
	linksChecker := service.NewHTTPLinkChecker(clientChecker)
	linkService := service.NewLinkService(linkRepo, linksChecker)
	pdfBuilder := pdf.NewPDFBuilder()
	var linkHandler linkHandlerI = handlers.NewLinkHandler(linkService, pdfBuilder)

	a.echo.POST("/links/get_status", linkHandler.GetStatus)
	a.echo.POST("/links/pdf", linkHandler.BuildPDF)

	// graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// start server in a separate goroutine
	go func() {
		if err := a.echo.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.echo.Logger.Fatal("shutting down the server unexpectedly:", err)
		}
	}()

	// wait termination signal
	<-ctx.Done()
	log.Println("Shutdown signal received. Starting graceful shutdown...")

	// Shutdown the Echo server with a timeout (Allows existing requests to finish)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// It's vital to wait for this 'Shutdown' call to complete
	err := a.echo.Shutdown(shutdownCtx)
	cancel() // Release the shutdown context resources immediately after the call

	if err != nil {
		a.echo.Logger.Fatal("Server forced to shutdown or timed out:", err)
	}
	log.Println("All active requests finished processing.")

	// Call the repository method to persist data (Task 2)
	log.Println("Persisting data to JSON files...")
	if err := dataManager.StoreDataToJSON(); err != nil {
		log.Printf("Error storing data to JSON: %v\n", err)
	} else {
		log.Println("Data persistence complete.")
	}

	log.Println("Server gracefully stopped.")
}
