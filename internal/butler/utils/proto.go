package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
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

func StructToMessage(obj interface{}, message proto.Message) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	err = jsonpb.UnmarshalString(string(b), message)
	if err != nil {
		return err
	}

	return nil
}
