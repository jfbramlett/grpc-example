package factories

import (
	"context"
	"reflect"
)

type InstanceCreator interface {
	NewInstance() interface{}
}

type reflectionInstanceCreator struct {
	typeOf		reflect.Type
}

func (i *reflectionInstanceCreator) NewInstance() interface{} {
	return reflect.New(i.typeOf).Interface()
}

func newReflectionInstanceCreator(ins interface{}) InstanceCreator {
	return &reflectionInstanceCreator{typeOf: reflect.TypeOf(ins)}
}


type contextInstanceCreator struct {}
func (c *contextInstanceCreator) NewInstance() interface{} {
	return context.Background()
}

func newContextInstanceCreator() InstanceCreator {
	return &contextInstanceCreator{}
}
