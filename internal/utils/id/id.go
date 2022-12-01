package id

import (
	b64 "encoding/base64"

	"github.com/octopipe/charlescd/internal/butler/errs"
)

func ToID(name string) string {
	return b64.StdEncoding.EncodeToString([]byte(name))
}

func DecodeID(id string) (string, error) {
	res, err := b64.StdEncoding.DecodeString(id)
	if err != nil {
		return "", errs.E(errs.Code("INVALID_ID"), err)
	}

	return string(res), nil
}
