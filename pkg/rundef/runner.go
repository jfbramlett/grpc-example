package rundef

import (
	"fmt"
	"github.com/jfbramlett/faker/fakegen"
	"github.com/jfbramlett/grpc-example/pkg/factories"
	valid "github.com/jfbramlett/grpc-example/pkg/validator"
	"log"
	"reflect"
	"testing"
)

type Runner interface {
	Run() RunResult
}

func NewRunner(runSuiteDef RunDefSuite, runDef RunDef, typeFactory factories.TypeFactory,
	clientFactory factories.ClientFactory, validator valid.Validator) Runner {
	return &basicRunner{runSuiteDef: runSuiteDef, runDef: runDef, typeFactory: typeFactory,
		clientFactory: clientFactory, validator: validator}
}

type basicRunner struct {
	runSuiteDef 		RunDefSuite
	runDef 				RunDef
	typeFactory			factories.TypeFactory
	clientFactory		factories.ClientFactory
	validator			valid.Validator
}

func (b *basicRunner) Run() RunResult {
	client, err := b.clientFactory.GetClient(b.runDef.ClientClassName)
	if err != nil {
		log.Println(err)
		return b.failedRun(err)
	}

	_, err =b.invokeTestAgainst(client)

	if err != nil {
		log.Println(fmt.Sprintf("Failed: %s", err))
		return b.failedRun(err)
	} else {
		log.Println("Success")
		return b.passedRun()
	}
}

func (b *basicRunner) invokeTestAgainst(any interface{}) (bool, error) {
	method, err := b.getMethod(any)
	if err != nil {
		return false, err
	}

	params, err := b.getParams(method)
	if err != nil {
		return false, err
	}

	result := method.Call(params)

	return b.validator.Validate(result)
}

func (b *basicRunner) getMethod(any interface{}) (reflect.Value, error) {
	meth := reflect.ValueOf(any).MethodByName(b.runDef.Function.Name)
	if meth.IsNil() {
		return reflect.ValueOf(""), fmt.Errorf("failed to find method %s", b.runDef.Function.Name)
	}

	return meth, nil
}

func (b *basicRunner) getParams(method reflect.Value) ([]reflect.Value, error) {
	methodType := method.Type()

	argCount := methodType.NumIn()
	if methodType.IsVariadic() {
		argCount--
	}

	argDef := b.runDef.Function.Args
	if argDef == nil {
		argDef = make(map[string]FunctionArg)
	}

	in := make([]reflect.Value, argCount)
	for i := 0; i < argCount; i++ {
		methodArgType := methodType.In(i)

		newInstance, err := b.createParam(methodArgType)
		if err != nil {
			return []reflect.Value{}, err
		}
		in[i] = reflect.ValueOf(newInstance)
	}

	return in, nil
}

func (b *basicRunner) failedRun(err error) RunResult {
	return RunResult{Name: b.runDef.Name, Passed: false, Error: err}
}

func (b *basicRunner) passedRun() RunResult {
	return RunResult{Name: b.runDef.Name, Passed: true}
}


func (f *basicRunner) createParam(argType reflect.Type) (interface{}, error) {
	insCreator, err := f.typeFactory.GetInstanceCreatorForType(argType)
	if err != nil {
		return nil, err
	}

	generator := f.getFakeGenerator()

	newInstance := insCreator.NewInstance()

	argName := factories.GetTypeName(argType)

	argDescription, found := f.runDef.Function.Args[argName]
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

func (f *basicRunner) getFakeGenerator() *fakegen.FakeGenerator {
	generator := fakegen.NewFakeGenerator()
	generator.AddFieldFilter("XXX_.*")
	if f.runSuiteDef.GlobalTags != nil {
		for k, v := range f.runSuiteDef.GlobalTags {
			generator.AddFieldTag(k, v)
		}
	}
	if f.runDef.RunTags != nil {
		for k, v := range f.runDef.RunTags {
			generator.AddFieldTag(k, v)
		}
	}
	if f.runSuiteDef.GlobalValues != nil {
		for k, v := range f.runSuiteDef.GlobalValues {
			generator.AddProvider(k, StaticTagProvider{val: v}.GetTaggedValue)
			generator.AddFieldTag(k, k)
		}
	}
	if f.runDef.RunValues != nil {
		for k, v := range f.runDef.RunValues {
			generator.AddProvider(k, StaticTagProvider{val: v}.GetTaggedValue)
			generator.AddFieldTag(k, k)
		}

	}
	return generator
}


type StaticTagProvider struct {
val			interface{}
}

func (s StaticTagProvider) GetTaggedValue(v reflect.Value) (interface{}, error) {
	return s.val, nil
}


// Variation of our runner that runs the test as a sub-test of the given test
type testingRunner struct {
	underlying		Runner
	name			string
	mainTest		*testing.T
}

func (b *testingRunner) Run() RunResult {
	var result RunResult
	b.mainTest.Run(b.name, func(t *testing.T){
		result = b.underlying.Run()
		if !result.Passed {
			t.Fail()
		}
	})
	return result
}

func NewTestingRunner(t *testing.T, name string, underlying Runner) Runner {
	return &testingRunner{underlying: underlying, name: name, mainTest: t}
}