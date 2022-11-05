package utils

import (
	"encoding/json"

	"github.com/gogo/protobuf/proto"
)

func MessageToStruct(message proto.Message, obj *interface{}) error {
	b, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	return nil
}

func StructToMessage(obj interface{}, message interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, message)
	if err != nil {
		return err
	}

	return nil
}
