package runner

import (
	"github.com/jfbramlett/grpc-example/pkg/factories"
	valid "github.com/jfbramlett/grpc-example/pkg/validator"
	"testing"
)

type RunnerFactory interface {
	GetRunner(runSuiteDef RunSuiteDef, runDef RunDef, typeFactory factories.TypeFactory, clientFactory factories.ClientFactory, validator valid.Validator) Runner
}

// runner factory that just creates a basic runner instance
type defaultRunnerFactory struct {}

func (d *defaultRunnerFactory) GetRunner(runSuiteDef RunSuiteDef, runDef RunDef, typeFactory factories.TypeFactory, clientFactory factories.ClientFactory, validator valid.Validator) Runner {
	return NewRunner(runSuiteDef, runDef, typeFactory, clientFactory, validator)
}

// runner factory that creates new runner instances wrapped in standard Go testing
type testingRunnerFactory struct {
	mainTest		*testing.T
}

func (d *testingRunnerFactory) GetRunner(runSuiteDef RunSuiteDef, runDef RunDef, typeFactory factories.TypeFactory, clientFactory factories.ClientFactory, validator valid.Validator) Runner {
	return NewTestingRunner(d.mainTest, runDef.Name, NewRunner(runSuiteDef, runDef, typeFactory, clientFactory, validator))
}

