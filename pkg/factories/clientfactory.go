package factories

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/routeguide"
	"google.golang.org/grpc"
	"log"
)

type ClientFactory interface {
	GetClient(name string) (interface{}, error)
	Close()
}

func GetClientFactory() ClientFactory {
	clientFactory := grpcClientFactory{}
	clientFactory.registerClients()
	return &clientFactory
}


type grpcClientFactory struct {
	registeredClients		map[string]interface{}
	grpcConnections			map[string]*grpc.ClientConn
}


func (r *grpcClientFactory) GetClient(name string) (interface{}, error) {
	client, found := r.registeredClients[name]
	if !found {
		err := fmt.Errorf("failed to find registered client with name %s", name)
		return nil, err
	}

	return client, nil
}

func (r *grpcClientFactory) Close() {
	for _, closeable := range r.grpcConnections {
		closeable.Close()
	}
}

func (r *grpcClientFactory) newGrpcClient(host string) *grpc.ClientConn {
	var conn *grpc.ClientConn
	var found bool
	if conn, found = r.grpcConnections[host]; !found {
		// prep our grpc env
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())

		dialConn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalln("fail to dial ", err)
			return nil
		}
		conn = dialConn
		r.grpcConnections[host] = conn
	}

	return conn

}
func (r *grpcClientFactory) registerClients() {
	r.registeredClients = make(map[string]interface{})
	r.grpcConnections = make(map[string]*grpc.ClientConn)

	r.registeredClients["routeguide.RouteGuideClient"] = routeguide.NewRouteGuideClient(r.newGrpcClient("localhost:2112"))
}