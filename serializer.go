package serializer

import (
	"errors"
	"reflect"
)

// Any is an Alias for interface{}
type Any = interface{}

// ISerializer provides a generic way to define Serializer
type ISerializer interface {
	Init(object Any) error
	Serialize() (map[string]Any, error)
}

// BaseSerializer implements a basic serializer to map[string]Any
type BaseSerializer struct {
	Object          Any
	cachedHash      map[string]func(*BaseSerializer) Any
	objectElemValue reflect.Value
}

// Init takes an object for serialization
func (ser *BaseSerializer) Init(object Any) error {
	if object == nil {
		return errors.New("object is nil")
	}
	ser.Object = object
	ser.objectElemValue = reflect.ValueOf(object).Elem()
	return nil
}

// RegisterFieldName defines a `name` field with its value from ser.Object[valueName]
func (ser *BaseSerializer) RegisterFieldName(name string, valueName string) {
	if field := ser.objectElemValue.FieldByName(valueName); field.CanInterface() {
		ser.RegisterFieldFunc(name, func(_ *BaseSerializer) Any {
			return field.Interface()
		})
	}
}

// RegisterFieldFunc defines a `name` field with the returned value from func.
func (ser *BaseSerializer) RegisterFieldFunc(name string, handler func(*BaseSerializer) Any) {
	ser.cachedHash[name] = handler
}

// Serialize generates the result.
func (ser *BaseSerializer) Serialize() map[string]Any {
	resultHash := make(map[string]Any)
	for k, v := range ser.cachedHash {
		resultHash[k] = v(ser)
	}
	return resultHash
}
