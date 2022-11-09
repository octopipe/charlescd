package errs

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHTTPResponse(c echo.Context, logger *zap.Logger, err error) error {
	const defaultMessage string = "internal server error - please contact support"
	var e *Error
	if errors.As(err, &e) {
		httpLog(logger, e.Kind, *e)
		switch e.Kind {
		case Internal, Database, Integration:
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
			return c.JSON(httpErrorStatusCode(e.Kind), err)
		}
	}

	logger.Error(err.Error(), zap.String("code", string(Unknown.String())))
	return c.JSON(http.StatusInternalServerError, HTTPError{
		Code:    Unknown.String(),
		Message: "Unknown Error",
	})
}

func httpLog(logger *zap.Logger, k Kind, e Error) {
	switch k {
	case NotExist:
		logger.Info(e.Error(), zap.String("code", string(e.Code)))
	default:
		logger.Error(e.Error(), zap.String("code", string(e.Code)))
	}
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
