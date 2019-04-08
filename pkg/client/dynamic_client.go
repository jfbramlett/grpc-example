package client

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/pkg/runner"
	"log"
)

func RunDynamicClient() {
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

