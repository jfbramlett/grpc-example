package server

import (
	"context"
	"fmt"
	"github.com/jfbramlett/faker/pkg/fakegen"
	"github.com/jfbramlett/grpc-example/interceptor"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/jfbramlett/grpc-example/routeguide"
)

type routeGuideServer struct {
}


func (s *routeGuideServer) FindRoute(ctx context.Context, r *routeguide.RouteRequest) (*routeguide.RouteDetails, error) {
	fmt.Println("Received request for FindRoute", r)
	fakegen.AddFieldFilter("XXX_.*")
	routeResponse := &routeguide.RouteDetails{}
	fakegen.FakeData(routeResponse)
	fmt.Println("Sending response: ", routeResponse)
	return routeResponse, nil
}


func newServer() *routeGuideServer {
	s := &routeGuideServer{}
	return s
}

func RunServer() {
	fmt.Println("Starting GRPC Server")
	lis, err := net.Listen("tcp", "localhost:2112")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{interceptor.WithServerUnaryInterceptor()}

	grpcServer := grpc.NewServer(opts...)
	routeguide.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}