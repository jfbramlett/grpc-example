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
		err := f.runTest(testdef, clientFactory, typeFactory)
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

func (f *basicTestSuite) runTest(test TestDef, clientFactory factories.ClientFactory, typeFactory factories.TypeFactory) error {
	client, err := clientFactory.GetClient(test.ClientClassName)
	if err != nil {
		log.Println(err)
		return err
	}

	response, err := f.invoke(test, client, test.Function.Name, typeFactory)

	if err != nil {
		log.Println(fmt.Sprintf("Failed: %s", err))
		return err
	} else {
		log.Println(fmt.Sprintf("Success: %v", response))
		return nil
	}

}

func (f *basicTestSuite) invoke(testDef TestDef, any interface{}, name string, factory factories.TypeFactory) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()

	argCount := methodType.NumIn()
	if methodType.IsVariadic() {
		argCount--
	}

	argDef := testDef.Function.Args
	if argDef == nil {
		argDef = make(map[string]FunctionArg)
	}

	in := make([]reflect.Value, argCount)
	for i := 0; i < argCount; i++ {
		methodArgType := methodType.In(i)

		newInstance, err := f.createParam(testDef.Function, methodArgType, factory)
		if err != nil {
			return reflect.ValueOf(""), err
		}
		in[i] = reflect.ValueOf(newInstance)
	}

	return method.Call(in)[0], nil
}

func (f *basicTestSuite) createParam(funcDef FunctionDef, argType reflect.Type, factory factories.TypeFactory) (interface{}, error) {
	insCreator, err := factory.GetInstanceCreatorForType(argType)
	if err != nil {
		return nil, err
	}

	newInstance := insCreator.NewInstance()

	argName := factories.GetTypeName(argType)

	argDescription, found := funcDef.Args[argName]
	if found && argDescription.FieldTags != nil {
		fakegen.SetFieldTags(argDescription.FieldTags)
	}
	defer fakegen.ClearFieldTags()
	fakegen.FakeData(newInstance)

	if found && argDescription.ValuesOverrideJson != "" {
		err := json.Unmarshal([]byte(argDescription.ValuesOverrideJson), newInstance)
		if err != nil {
			return nil, err
		}
	}
	return newInstance, nil
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


func NewAutoTestSuite(interfaceType reflect.Type) (TestSuite, error) {
	testSuite := TestSuiteDef{Tests: make([]TestDef, 0)}

	for i := 0; i < interfaceType.NumMethod(); i++ {
		methodName := interfaceType.Method(i).Name
		testSuite.Tests = append(testSuite.Tests, TestDef{Name: fmt.Sprintf(" Test %s.%s", factories.GetTypeName(interfaceType), methodName),
			ClientClassName: factories.GetTypeName(interfaceType),
			Function: FunctionDef{Name: methodName},
		})

	}

	return &basicTestSuite{testSuiteDef: testSuite}, nil
}

