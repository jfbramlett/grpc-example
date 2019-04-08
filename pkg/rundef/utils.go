package rundef

import (
	"encoding/json"
	"fmt"
	"github.com/jfbramlett/grpc-example/pkg/factories"
	valid "github.com/jfbramlett/grpc-example/pkg/validator"
	"io/ioutil"
	"log"
	"reflect"
)

// function used to build a run suite from a file, the file is a JSON file containing the definition of what to run
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

// function used to build a run suite given a type - this uses reflection to identify the methods to wrap
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


func getValidator(runDef RunDef, validatorFactory valid.ValidatorFactory) (valid.Validator, error) {
	if runDef.Validator == "" {
		return &valid.DefaultValidator{}, nil
	} else {
		validator, err := validatorFactory.GetValidator(runDef.Validator)
		if err != nil {
			return nil, err
		}
		return validator, nil
	}
}

