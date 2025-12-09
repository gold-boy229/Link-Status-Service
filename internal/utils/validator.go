package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type customValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *customValidator {
	return &customValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

func (cv *customValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Return the error so Echo's error handler can process it
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
