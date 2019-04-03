package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jfbramlett/faker/pkg/fakegen"
	"github.com/jfbramlett/grpc-example/pkg/funcdef"
	"github.com/jfbramlett/grpc-example/pkg/typefactory"
	"github.com/jfbramlett/grpc-example/routeguide"
	"google.golang.org/grpc"
	"log"
	"reflect"
)


func RunClient() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:2112", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := routeguide.NewRouteGuideClient(conn)

	fakegen.AddFieldFilter("XXX_.*")
	typeFactory := typefactory.GetTypeFactory()

	testDef := funcdef.Function{Name: "FindRoute", Args: []funcdef.FunctionArg {{FuncType: "context.Context"}, {FuncType: "routeguide.RouteRequest", ValuesOverride: "{\"userid\": 15}", FieldTags: map[string]string {"Email": "email"}}}}

	args, err := prepareArgs(testDef.Args, typeFactory)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed: %s", err))
		return
	}

	// Looking for a valid feature
	response, err := invoke(client, testDef.Name, args...)

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed: %s", err))
	} else {
		fmt.Println(fmt.Sprintf("Success: %v", response))
	}
}


func prepareArgs(funcArgs []funcdef.FunctionArg, factory typefactory.TypeFactory) ([]interface{}, error) {
	args := make([]interface{}, len(funcArgs))

	for i := 0; i < len(funcArgs); i++ {
		fa := funcArgs[i]

		insCreator, err := factory.GetInstanceCreator(fa.FuncType)
		if err != nil {
			return nil, err
		}

		newInstance := insCreator.NewInstance()
		if len(fa.FieldTags) > 0 {
			fakegen.SetFieldTags(fa.FieldTags)
		}
		defer fakegen.ClearFieldTags()
		fakegen.FakeData(newInstance)


		if fa.ValuesOverride != "" {
			err := json.Unmarshal([]byte(fa.ValuesOverride), newInstance)
			if err != nil {
				return nil, err
			}
		}
		args[i] = newInstance
	}

	return args, nil
}


func hasOverride(fieldname string, overrides map[string]string) bool {
	_, found := overrides[fieldname]
	return found
}


func invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()
	if methodType.IsVariadic() && numIn < len(args) - 1 {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var methArgType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			methArgType = methodType.In(numIn - 1).Elem()
		} else {
			methArgType = methodType.In(i)
		}

		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, argValue.String(), methArgType)
		}

		argValueType := argValue.Type()
		if argValueType.Kind() == reflect.Ptr && methArgType.Kind() == reflect.Ptr {
			in[i] = argValue
		} else if argValueType.ConvertibleTo(methArgType) {
			in[i] = argValue.Convert(methArgType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, argValueType, methArgType)
		}
	}
	return method.Call(in)[0], nil
}