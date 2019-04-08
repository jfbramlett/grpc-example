package client

import (
	"fmt"
	"github.com/jfbramlett/go-dynamic-runner/factories"
	"github.com/jfbramlett/go-dynamic-runner/runner"
	"github.com/jfbramlett/grpc-example/routeguide"
	"google.golang.org/grpc"
	"log"
)

func RunDynamicClient() {
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


	suite, err := runner.NewRunSuite("testdata/runsuite.json")
	/*
	suite, err := rundef.NewAutoRunSuite(reflect.TypeOf((*routeguide.RouteGuideClient)(nil)).Elem(),
		map[string]interface{} {"Destination": "UNC"},
		map[string]string {"Email": "email"},
		[]string {})
	 */
	if err != nil {
		log.Fatalln(err)
		return
	}

	results := suite.Run()
	passedCount := 0
	failedCount := 0
	for _, r := range results {
		if r.Passed {
			log.Println(fmt.Sprintf("Test %s - PASSED", r.Name))
			passedCount++
		} else {
			log.Println(fmt.Sprintf("Test %s - FAILED", r.Name))
			failedCount++
		}
	}

	log.Println(fmt.Sprintf("Total Tests: %d", passedCount + failedCount))
	log.Println(fmt.Sprintf("Passed Tests: %d", passedCount))
	log.Println(fmt.Sprintf("Failed Tests: %d", failedCount))


}

