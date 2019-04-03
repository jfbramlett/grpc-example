package factories

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/routeguide"
	"reflect"
)

type TypeFactory interface {
	GetInstanceCreator(name string) (InstanceCreator, error)
}

func GetTypeFactory() TypeFactory {
	typeFactory := &grpcTypeFactory{}
	typeFactory.registerTypes()
	return typeFactory
}

func typeNameFromIns(ins interface{}) string {
	return fmt.Sprintf("%s", reflect.TypeOf(ins))
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
	g.typeMap = make(map[string]InstanceCreator)

	g.typeMap["context.Context"] = newContextInstanceCreator()
	g.typeMap[typeNameFromIns(routeguide.RouteRequest{})] = newReflectionInstanceCreator(routeguide.RouteRequest{})
	g.typeMap[typeNameFromIns(routeguide.RouteDetails{})] = newReflectionInstanceCreator(routeguide.RouteDetails{})
}

