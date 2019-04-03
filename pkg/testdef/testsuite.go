package testdef

import (
	"encoding/json"
	"fmt"
	"github.com/jfbramlett/faker/pkg/fakegen"
	"github.com/jfbramlett/grpc-example/pkg/factories"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"reflect"
)

type TestSuite interface {
	Run() []TestResult
}

type TestResult struct {
	Name		string
	Passed		bool
	Error		error
}


type basicTestSuite struct {
	testSuiteDef		TestSuiteDef
}

func (f *basicTestSuite) Run() []TestResult {
	testResults := make([]TestResult, 0)

	// prep our grpc env
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:2112", opts...)
	if err != nil {
		log.Fatalln("fail to dial ", err)
		return testResults
	}
	defer conn.Close()

	fakegen.AddFieldFilter("XXX_.*")
	typeFactory := factories.GetTypeFactory()
	clientFactory := factories.GetClientFactory(conn)


	for _, testdef := range f.testSuiteDef.Tests {
		err := runTest(testdef, clientFactory, typeFactory)
		if err != nil {
			log.Println(fmt.Sprintf("Test: %s - FAILED   %s", testdef.Name, err))
			testResults = append(testResults, TestResult{Name: testdef.Name, Passed: false, Error: err})
		} else {
			log.Println(fmt.Sprintf("Test: %s - PASSED", testdef.Name))
			testResults = append(testResults, TestResult{Name: testdef.Name, Passed: true})
		}
	}

	return testResults
}

func runTest(test TestDef, clientFactory factories.ClientFactory, typeFactory factories.TypeFactory) error {
	client, err := clientFactory.GetClient(test.ClientClassName)
	if err != nil {
		log.Println(err)
		return err
	}

	args, err := prepareArgs(test.Function.Args, typeFactory)
	if err != nil {
		log.Println(fmt.Sprintf("Failed: %s", err))
		return err
	}

	// Looking for a valid feature
	response, err := invoke(client, test.Function.Name, args...)

	if err != nil {
		log.Println(fmt.Sprintf("Failed: %s", err))
		return err
	} else {
		log.Println(fmt.Sprintf("Success: %v", response))
		return nil
	}

}


func prepareArgs(funcArgs []FunctionArg, factory factories.TypeFactory) ([]interface{}, error) {
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


		if fa.ValuesOverrideJson != "" {
			err := json.Unmarshal([]byte(fa.ValuesOverrideJson), newInstance)
			if err != nil {
				return nil, err
			}
		}
		args[i] = newInstance
	}

	return args, nil
}

func invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()
	if methodType.IsVariadic() && numIn < len(args) - 1 {
		return reflect.ValueOf(nil), fmt.Errorf("method %s must have minimum %d params have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("method %s must have %d params have %d", name, numIn, len(args))
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
			return reflect.ValueOf(nil), fmt.Errorf("method %s param[%d] must be %s have %s", name, i, argValue.String(), methArgType)
		}

		argValueType := argValue.Type()
		if argValueType.Kind() == reflect.Ptr && methArgType.Kind() == reflect.Ptr {
			in[i] = argValue
		} else if argValueType.ConvertibleTo(methArgType) {
			in[i] = argValue.Convert(methArgType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("method %s param[%d] must be %s have %s", name, i, argValueType, methArgType)
		}
	}
	return method.Call(in)[0], nil
}

func NewTestSuite(configFile string) (TestSuite, error) {
	testSuite := TestSuiteDef{}
	testSuiteDef, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	err = json.Unmarshal(testSuiteDef, &testSuite)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &basicTestSuite{testSuiteDef: testSuite}, nil
}


