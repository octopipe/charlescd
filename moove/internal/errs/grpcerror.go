package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParseGrpcError(err error) error {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			return E(NotExist, err)
		default:
			return E(Integration, err)
		}
	}

	return E(Unknown, err)
}
