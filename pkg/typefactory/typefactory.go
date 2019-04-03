package typefactory

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/routeguide"
	"reflect"
)

type TypeFactory interface {
	GetInstanceCreator(name string) (InstanceCreator, error)
}


type grpcTypeFactory struct {
	typeMap 		map[string]InstanceCreator
}


func (g *grpcTypeFactory) GetInstanceCreator(name string) (InstanceCreator, error) {
	t, found := g.typeMap[name]

	if !found {
		return nil, fmt.Errorf("no type named %s registered", name)
	} else {
		return t, nil
	}
}

func (g *grpcTypeFactory) registerTypes() {
	g.typeMap["context.Context"] = newContextInstanceCreator()
	g.typeMap[typeNameFromIns(routeguide.RouteRequest{})] = newReflectionInstanceCreator(routeguide.RouteRequest{})
	g.typeMap[typeNameFromIns(routeguide.RouteDetails{})] = newReflectionInstanceCreator(routeguide.RouteDetails{})
}

func typeNameFromIns(ins interface{}) string {
	return fmt.Sprintf("%s", reflect.TypeOf(ins))
}


func GetTypeFactory() TypeFactory {
	typeFactory := &grpcTypeFactory{typeMap: make(map[string]InstanceCreator)}
	typeFactory.registerTypes()
	return typeFactory
}
