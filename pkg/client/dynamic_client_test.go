package client

import (
	"github.com/jfbramlett/grpc-example/pkg/rundef"
	"github.com/jfbramlett/grpc-example/routeguide"
	"reflect"
	"testing"
)

func TestRunDynamicClient(t *testing.T) {
	runner, err := rundef.NewAutoTestingRunSuite(t,
		reflect.TypeOf((*routeguide.RouteGuideClient)(nil)).Elem(),
		map[string]interface{} {"Destination": "UNC"},
		map[string]string {"Email": "email"},
		[]string {})

	if err != nil {
		t.Fail()
	}

	runner.Run()
}
