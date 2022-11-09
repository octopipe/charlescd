package errs

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

type Kind uint8

type Code string

type ErrorField struct {
	name       string
	validation string
}

type Error struct {
	Kind   Kind
	Code   Code
	Err    error
	Fields []ErrorField
}

const (
	Other         Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                   // Invalid operation for this type of item.
	IO                        // External I/O error such as network failure.
	Exist                     // Item already exists.
	NotExist                  // Item does not exist.
	Private                   // Information withheld.
	Internal                  // Internal error or inconsistency.
	BrokenLink                // Link target does not exist.
	Database                  // Error from database.
	Validation                // Input validation error.
	Unanticipated             // Unanticipated error.
	Unknown
	InvalidRequest // Invalid Request
	Unauthenticated
	Unauthorized
	Integration
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other_error"
	case Invalid:
		return "invalid_operation"
	case IO:
		return "I/O_error"
	case Exist:
		return "item_already_exists"
	case NotExist:
		return "item_does_not_exist"
	case BrokenLink:
		return "link_target_does_not_exist"
	case Private:
		return "information_withheld"
	case Internal:
		return "internal_error"
	case Database:
		return "database_error"
	case Validation:
		return "input_validation_error"
	case Unanticipated:
		return "unanticipated_error"
	case Unknown:
		return "unknown_error"
	case InvalidRequest:
		return "invalid_request_error"
	case Unauthenticated:
		return "unauthenticated_request"
	case Unauthorized:
		return "unauthorized_request"
	case Integration:
		return "integration_request"
	}
	return "unknown_error_kind"
}

func (e Error) Error() string {
	return e.Err.Error()
}

func E(args ...interface{}) error {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Err = errors.New(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			e.Err = arg
		case error:
			// if the error implements stackTracer, then it is
			// a pkg/errors error type and does not need to have
			// the stack added
			_, ok := arg.(stackTracer)
			if ok {
				e.Err = arg
			} else {
				e.Err = errors.WithStack(arg)
			}
		case Code:
			e.Code = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			return fmt.Errorf("errors.E: bad call from %s:%d: %v, unknown type %T, value %v in error call", file, line, args, arg, arg)
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}
	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	if prev.Code == e.Code {
		prev.Code = ""
	}
	// If this error has Code == "", pull up the inner one.
	if e.Code == "" {
		e.Code = prev.Code
		prev.Code = ""
	}

	return e
}
