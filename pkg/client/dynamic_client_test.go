package client

import (
	"github.com/jfbramlett/go-dynamic-runner/factories"
	"github.com/jfbramlett/go-dynamic-runner/runner"
	"github.com/jfbramlett/grpc-example/routeguide"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"testing"
)

func TestRunDynamicClient(t *testing.T) {
	connections := make([]*grpc.ClientConn, 0)
	defer func() {
		for _, c := range connections {
			c.Close()
		}
	}()

	conn := newGrpcClient("localhost:2112")
	if conn == nil {
		log.Fatalln("failed to resolve grpc connection")
	}
	connections = append(connections, conn)
	factories.GlobalClientFactory.RegisterClient("routeguide.RouteGuideClient", routeguide.NewRouteGuideClient(conn))

	factories.GlobalTypeFactory.RegisterType(factories.GetTypeNameFromIns(routeguide.RouteRequest{}), factories.NewReflectionInstanceCreator(routeguide.RouteRequest{}))
	factories.GlobalTypeFactory.RegisterType(factories.GetTypeNameFromIns(routeguide.RouteDetails{}), factories.NewReflectionInstanceCreator(routeguide.RouteDetails{}))

	runner, err := runner.NewAutoTestingRunSuite(t,
		reflect.TypeOf((*routeguide.RouteGuideClient)(nil)).Elem(),
		map[string]interface{} {"Destination": "UNC"},
		map[string]string {"Email": "email"},
		[]string {})

	if err != nil {
		t.Fail()
	}

	runner.Run()
}
