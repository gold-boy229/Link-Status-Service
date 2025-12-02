package app

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Return the error so Echo's error handler can process it
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewCustomValidator() *customValidator {
	return &customValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
}
