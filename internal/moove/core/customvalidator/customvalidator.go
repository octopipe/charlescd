package customvalidator

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

type ValidateError struct {
	Field string
	Tag   string
	Value string
}

func NewCustomValidator() CustomValidator {
	return CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	errors := []ValidateError{}
	if err := cv.validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			newValidateError := ValidateError{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			}
			errors = append(errors, newValidateError)
		}
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}
	return nil
}
