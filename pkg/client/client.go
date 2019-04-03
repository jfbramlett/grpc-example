package client

import (
	"context"
	"flag"
	"github.com/jfbramlett/grpc-example/routeguide"
	"google.golang.org/grpc"
	"log"
)


func RunClient() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:2112", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := routeguide.NewRouteGuideClient(conn)

	request := &routeguide.RouteRequest{Destination: "NC", Email: "something@gmail.com", Userid: 10}

	response, err := client.FindRoute(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Success: ", response)
}

