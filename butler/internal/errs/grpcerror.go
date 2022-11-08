package errs

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCError(logger logr.Logger, err error) error {

	var e *Error
	if errors.As(err, &e) {
		st := status.New(grpcErrorCode(e.Kind), e.Err.Error())

		if e.Fields != nil {
			br := &errdetails.BadRequest{}
			for _, field := range e.Fields {
				v := &errdetails.BadRequest_FieldViolation{
					Field:       field.name,
					Description: fmt.Sprintf("field %s is %s", field.name, field.validation),
				}

				br.FieldViolations = append(br.FieldViolations, v)
			}
		}

		grpcLog(logger, e.Kind, e)
		return st.Err()
	}

	logger.Error(err, err.Error())
	return status.New(codes.Unknown, "unknown_error").Err()
}

func grpcLog(logger logr.Logger, kind Kind, e *Error) {
	switch kind {
	case NotExist:
		logger.Info(e.Error(), "code", e.Code)
		return
	default:
		logger.Error(e, e.Error(), "code", e.Code, "fields", e.Fields)
		return
	}
}

func grpcErrorCode(kind Kind) codes.Code {
	switch kind {
	default:
		return codes.Unknown
	}
}
