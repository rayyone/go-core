package maps

import (
	"encoding/json"
	"reflect"

	"github.com/rayyone/go-core/errors"
)

// Sanitize Remove key with nil value
func Sanitize(m map[string]interface{}) error {
	for key, val := range m {
		if reflect.ValueOf(val).IsNil() {
			delete(m, key)
		}
	}

	return nil
}

// ConvertToStruct Decode map[string]interface{} to struct
func ConvertToStruct(input interface{}, output interface{}) error {
	i, ok := input.(map[string]interface{})
	if !ok {
		return errors.New("Cannot convert maps to struct, input is not a type of map[string]interface{}")
	}

	inputBs, _ := json.Marshal(i)
	if err := json.Unmarshal(inputBs, output); err != nil {
		return errors.BadRequest.Newf("Cannot unmarshal when converting to struct. Error: %v", err)
	}

	return nil
}
