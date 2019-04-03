package factories

import (
	"fmt"
	"github.com/jfbramlett/grpc-example/routeguide"
	"google.golang.org/grpc"
)

type ClientFactory interface {
	GetClient(name string) (interface{}, error)
}

func GetClientFactory(conn *grpc.ClientConn) ClientFactory {
	clientFactory := simpleClientFactory{}
	clientFactory.registerClients(conn)
	return &clientFactory
}


type simpleClientFactory struct {
	registeredClients		map[string]interface{}
}


func (r *simpleClientFactory) GetClient(name string) (interface{}, error) {
	client, found := r.registeredClients[name]
	if !found {
		err := fmt.Errorf("failed to find registered client with name %s", name)
		return nil, err
	}

	return client, nil
}



func (r *simpleClientFactory) registerClients(conn *grpc.ClientConn) {
	r.registeredClients = make(map[string]interface{})

	r.registeredClients["routeguide.RouteGuideClient"] = routeguide.NewRouteGuideClient(conn)
}