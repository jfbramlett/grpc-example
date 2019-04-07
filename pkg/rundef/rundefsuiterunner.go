package rundef

import (
	"encoding/json"
	"fmt"
	"github.com/jfbramlett/grpc-example/pkg/factories"
	"io/ioutil"
	"log"
	"reflect"
)

type RunSuiteRunner interface {
	Run() []RunResult
}

type RunResult struct {
	Name		string
	Result		reflect.Value
	Passed		bool
	Error		error
}


type basicRunSuite struct {
	runSuite		RunDefSuite
}

func (f *basicRunSuite) Run() []RunResult {
	testResults := make([]RunResult, 0)

	clientFactory := factories.GetClientFactory()
	typeFactory := factories.GetTypeFactory()

	for _, runDef := range f.runSuite.Tests {
		runner := NewRunDefRunner(f.runSuite, runDef, typeFactory, clientFactory)
		result := runner.Run()
		testResults = append(testResults, result)
	}

	clientFactory.Close()
	typeFactory.Close()

	return testResults
}

func NewRunSuite(configFile string) (RunSuiteRunner, error) {
	runSuite, err := buildRunSuiteFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return &basicRunSuite{runSuite: runSuite}, nil
}


func NewAutoRunSuite(interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string) (RunSuiteRunner, error) {
	runSuite, err := buildRunSuiteFromType(interfaceType, globalValues, globalTags, excludes)
	if err != nil {
		return nil, err
	}

	return &basicRunSuite{runSuite: runSuite}, nil
}



func buildRunSuiteFromFile(configFile string) (RunDefSuite, error) {
	runSuite := RunDefSuite{}
	runSuiteDef, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
		return RunDefSuite{}, err
	}

	err = json.Unmarshal(runSuiteDef, &runSuite)
	if err != nil {
		log.Fatalln(err)
		return RunDefSuite{}, err
	}

	return runSuite, nil
}

func buildRunSuiteFromType(interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string) (RunDefSuite, error) {
	runSuite := RunDefSuite{Tests: make([]RunDef, 0), GlobalValues: globalValues, GlobalTags: globalTags}

	for i := 0; i < interfaceType.NumMethod(); i++ {
		methodName := interfaceType.Method(i).Name
		for _, excludedMethod := range excludes {
			if methodName == excludedMethod {
				continue
			}
		}
		runSuite.Tests = append(runSuite.Tests, RunDef{Name: fmt.Sprintf(" Test %s.%s", factories.GetTypeName(interfaceType), methodName),
			ClientClassName: factories.GetTypeName(interfaceType),
			Function: FunctionDef{Name: methodName},
		})

	}
	return runSuite, nil
}