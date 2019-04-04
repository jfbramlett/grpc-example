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

	typeFactory := factories.GetTypeFactory()
	clientFactory := factories.GetClientFactory(conn)


	for _, testdef := range f.testSuiteDef.Tests {
		err := f.runTest(f.testSuiteDef, testdef, clientFactory, typeFactory)
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

func (f *basicTestSuite) runTest(testSuiteDef TestSuiteDef, test TestDef, clientFactory factories.ClientFactory, typeFactory factories.TypeFactory) error {
	client, err := clientFactory.GetClient(test.ClientClassName)
	if err != nil {
		log.Println(err)
		return err
	}

	response, err := f.invoke(testSuiteDef, test, client, test.Function.Name, typeFactory)

	if err != nil {
		log.Println(fmt.Sprintf("Failed: %s", err))
		return err
	} else {
		log.Println(fmt.Sprintf("Success: %v", response))
		return nil
	}

}

func (f *basicTestSuite) invoke(testSuiteDef TestSuiteDef, testDef TestDef, any interface{}, name string, factory factories.TypeFactory) (reflect.Value, error) {
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

		newInstance, err := f.createParam(testSuiteDef, testDef, testDef.Function, methodArgType, factory)
		if err != nil {
			return reflect.ValueOf(""), err
		}
		in[i] = reflect.ValueOf(newInstance)
	}

	return method.Call(in)[0], nil
}

func (f *basicTestSuite) createParam(testSuiteDef TestSuiteDef, testDef TestDef, funcDef FunctionDef, argType reflect.Type, factory factories.TypeFactory) (interface{}, error) {
	insCreator, err := factory.GetInstanceCreatorForType(argType)
	if err != nil {
		return nil, err
	}

	generator := f.getFakeGenerator(testSuiteDef, testDef)

	newInstance := insCreator.NewInstance()

	argName := factories.GetTypeName(argType)

	argDescription, found := funcDef.Args[argName]
	if found && argDescription.FieldTags != nil {
		for k, v := range argDescription.FieldTags {
			generator.AddFieldTag(k, v)
		}
	}

	if found && argDescription.ValuesOverride != nil {
		for k, v := range argDescription.ValuesOverride {
			generator.AddProvider(k, StaticTagProvider{val: v}.GetTaggedValue)
			generator.AddFieldTag(k, k)
		}
	}

	generator.FakeData(newInstance)

	return newInstance, nil
}

func (f *basicTestSuite) getFakeGenerator(testSuiteDef TestSuiteDef, testDef TestDef) *fakegen.FakeGenerator {
	generator := fakegen.NewFakeGenerator()
	generator.AddFieldFilter("XXX_.*")
	if testSuiteDef.GlobalTags != nil {
		for k, v := range testSuiteDef.GlobalTags {
			generator.AddFieldTag(k, v)
		}
	}
	if testDef.TestTags != nil {
		for k, v := range testDef.TestTags {
			generator.AddFieldTag(k, v)
		}
	}
	if testSuiteDef.GlobalValues != nil {
		for k, v := range testSuiteDef.GlobalValues {
			generator.AddProvider(k, StaticTagProvider{val: v}.GetTaggedValue)
			generator.AddFieldTag(k, k)
		}
	}
	if testDef.TestValues != nil {
		for k, v := range testDef.TestValues {
			generator.AddProvider(k, StaticTagProvider{val: v}.GetTaggedValue)
			generator.AddFieldTag(k, k)
		}

	}
	return generator
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


func NewAutoTestSuite(interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string) (TestSuite, error) {
	testSuite := TestSuiteDef{Tests: make([]TestDef, 0), GlobalValues: globalValues, GlobalTags: globalTags}

	for i := 0; i < interfaceType.NumMethod(); i++ {
		methodName := interfaceType.Method(i).Name
		for _, excludedMethod := range excludes {
			if methodName == excludedMethod {
				continue
			}
		}
		testSuite.Tests = append(testSuite.Tests, TestDef{Name: fmt.Sprintf(" Test %s.%s", factories.GetTypeName(interfaceType), methodName),
			ClientClassName: factories.GetTypeName(interfaceType),
			Function: FunctionDef{Name: methodName},
		})

	}

	return &basicTestSuite{testSuiteDef: testSuite}, nil
}



type StaticTagProvider struct {
	val			interface{}
}

func (s StaticTagProvider) GetTaggedValue(v reflect.Value) (interface{}, error) {
	return s.val, nil
}