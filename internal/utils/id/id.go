package id

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/octopipe/charlescd/internal/moove/errs"
)

func ToID(name string) string {
	return b64.RawURLEncoding.EncodeToString([]byte(name))
}

func DecodeID(id string) (string, error) {
	res, err := b64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return "", errs.E(errs.Invalid, errs.Code("INVALID_ID"), fmt.Errorf("invalid id %s", id))
	}

	return string(res), nil
}
