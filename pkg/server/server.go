package server

import (
	"context"
	"fmt"
	"github.com/jfbramlett/faker/fakegen"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/jfbramlett/grpc-example/routeguide"
)

type routeGuideServer struct {
	generator		*fakegen.FakeGenerator
}


func (s *routeGuideServer) FindRoute(ctx context.Context, r *routeguide.RouteRequest) (*routeguide.RouteDetails, error) {
	fmt.Println("Received request for FindRoute", r)

	s.generator.AddFieldFilter("XXX_.*")
	routeResponse := &routeguide.RouteDetails{}
	s.generator.FakeData(routeResponse)
	fmt.Println("Sending response: ", routeResponse)
	return routeResponse, nil
}


func newServer() *routeGuideServer {
	s := &routeGuideServer{generator: fakegen.NewFakeGenerator()}
	return s
}

func RunServer() {
	fmt.Println("Starting GRPC Server")
	lis, err := net.Listen("tcp", "localhost:2112")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{WithServerUnaryInterceptor()}

	grpcServer := grpc.NewServer(opts...)
	routeguide.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}