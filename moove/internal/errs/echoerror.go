package errs

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHTTPResponse(c echo.Context, err error) error {
	const defaultMessage string = "internal server error - please contact support"
	var e *Error
	if errors.As(err, &e) {
		switch e.Kind {
		case Internal, Database:
			err := HTTPError{
				Code:    e.Kind.String(),
				Message: defaultMessage,
			}
			return c.JSON(httpErrorStatusCode(e.Kind), err)
		default:
			err := HTTPError{
				Code:    e.Kind.String(),
				Message: e.Err.Error(),
			}
			c.JSON(httpErrorStatusCode(e.Kind), err)
		}
	}

	return c.JSON(http.StatusInternalServerError, HTTPError{
		Code:    Unknown.String(),
		Message: "Unknown Error",
	})
}

func httpErrorStatusCode(k Kind) int {
	switch k {
	case NotExist:
		return http.StatusNotFound
	case Invalid, Exist, Private, BrokenLink, Validation, InvalidRequest:
		return http.StatusBadRequest
	case Other, IO, Internal, Database, Unanticipated:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
