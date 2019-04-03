package client

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/pkg/testdef"
	"github.com/jfbramlett/grpc-example/routeguide"
	"log"
	"reflect"
)

func RunDynamicClient() {
	//suite, err := testdef.NewTestSuite("testdata/testsuite.json")
	suite, err := testdef.NewAutoTestSuite(reflect.TypeOf((*routeguide.RouteGuideClient)(nil)).Elem())
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

