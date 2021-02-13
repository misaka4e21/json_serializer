package serializer

import (
	"errors"
	"reflect"
)

// Any is an Alias for interface{}
type Any = interface{}

// ISerializer provides a generic way to define Serializer
type ISerializer interface {
	Init(object Any) ISerializer
	Serialize() (map[string]Any, error)
	SerializeIgnoreNull() (map[string]Any, error)
}

// BaseSerializer implements a basic serializer to map[string]Any
type BaseSerializer struct {
	Object          Any
	cachedHash      map[string]func(*BaseSerializer) Any
	objectElemValue reflect.Value
	serializeError  error
}

// Init takes an object for serialization
func (ser *BaseSerializer) Init(object Any) ISerializer {
	if object == nil {
		ser.serializeError = errors.New("object is nil")
	}
	ser.Object = object
	ser.objectElemValue = reflect.ValueOf(object).Elem()
	ser.cachedHash = make(map[string]func(*BaseSerializer) Any)
	return ser
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
func (ser *BaseSerializer) Serialize() (map[string]Any, error) {
	resultHash := make(map[string]Any)
	for k, v := range ser.cachedHash {
		resultHash[k] = v(ser)
	}
	return resultHash, ser.serializeError
}

// SerializeIgnoreNull generates the result, with nils ignored.
func (ser *BaseSerializer) SerializeIgnoreNull() (map[string]Any, error) {
	resultHash := make(map[string]Any)
	for k, v := range ser.cachedHash {
		value := v(ser)
		if value != nil {
			resultHash[k] = value
		}
	}
	return resultHash, ser.serializeError
}
