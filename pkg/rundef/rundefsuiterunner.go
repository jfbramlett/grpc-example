package rundef

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/pkg/factories"
	valid "github.com/jfbramlett/grpc-example/pkg/validator"
	"reflect"
	"testing"
)

type RunSuiteRunner interface {
	Run() []RunResult
}

// simple runner that just runs the tests and reports the results
type basicRunSuite struct {
	runSuite		RunDefSuite
}

func (f *basicRunSuite) Run() []RunResult {
	testResults := make([]RunResult, 0)

	clientFactory := factories.GetClientFactory()
	typeFactory := factories.GetTypeFactory()
	validatorFactory := valid.GetValidatorFactory()

	for _, runDef := range f.runSuite.Tests {
		validator, err := getValidator(runDef, validatorFactory)
		if err != nil {
			testResults = append(testResults, RunResult{Name: runDef.Name, Passed: false, Error: fmt.Errorf("failed to find configured validator %s", runDef.Validator)})
			continue
		}
		runner := NewRunDefRunner(f.runSuite, runDef, typeFactory, clientFactory, validator)
		result := runner.Run()
		testResults = append(testResults, result)
	}

	clientFactory.Close()
	typeFactory.Close()

	return testResults
}

// version of our RunSuite that wraps the tests as testing.T methods
type testingRunSuite struct {
	runSuite		RunDefSuite
	mainTest		*testing.T
}

func (r *testingRunSuite) Run() []RunResult {
	testResults := make([]RunResult, 0)

	clientFactory := factories.GetClientFactory()
	typeFactory := factories.GetTypeFactory()
	validatorFactory := valid.GetValidatorFactory()

	for _, runDef := range r.runSuite.Tests {
		validator, err := getValidator(runDef, validatorFactory)
		if err != nil {
			testResults = append(testResults, RunResult{Name: runDef.Name, Passed: false, Error: fmt.Errorf("failed to find configured validator %s", runDef.Validator)})
		}

		r.mainTest.Run(runDef.Name, func(t *testing.T){
			runner := NewRunDefRunner(r.runSuite, runDef, typeFactory, clientFactory, validator)
			result := runner.Run()
			testResults = append(testResults, result)
			if !result.Passed {
				t.Fail()
			}
		})
	}

	clientFactory.Close()
	typeFactory.Close()

	return testResults
}



// constructor to build a new RunSuite (set of things to execute). This builds it from a JSON-based config
func NewRunSuite(configFile string) (RunSuiteRunner, error) {
	runSuite, err := buildRunSuiteFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return &basicRunSuite{runSuite: runSuite}, nil
}

// constructor to build a new RunSuite (set of things to execute). This builds it from a JSON-based config. Each RunDef
// when executing will be run as a Go Test
func NewTestingRunSuite(t *testing.T, configFile string) (RunSuiteRunner, error) {
	runSuite, err := buildRunSuiteFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return &testingRunSuite{runSuite: runSuite, mainTest: t}, nil
}

// constructor method for creating a new RunSuite, a run suite represents a set of run defs (or things to run), this builds
// the suite automatically based on the type
func NewAutoRunSuite(interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string) (RunSuiteRunner, error) {
	runSuite, err := buildRunSuiteFromType(interfaceType, globalValues, globalTags, excludes)
	if err != nil {
		return nil, err
	}

	return &basicRunSuite{runSuite: runSuite}, nil
}

// constructor method for creating a new RunSuite, a run suite represents a set of run defs (or things to run), this builds
// the suite automatically based on the type. Each execution will be wrapped as a Go test case.
func NewAutoTestingRunSuite(t *testing.T, interfaceType reflect.Type, globalValues map[string]interface{}, globalTags map[string]string, excludes []string) (RunSuiteRunner, error) {
	runSuite, err := buildRunSuiteFromType(interfaceType, globalValues, globalTags, excludes)
	if err != nil {
		return nil, err
	}

	return &testingRunSuite{runSuite: runSuite, mainTest: t}, nil
}
