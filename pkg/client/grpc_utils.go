package client

import (
	"google.golang.org/grpc"
	"log"
)

func  newGrpcClient(host string) *grpc.ClientConn {
	// prep our grpc env
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Fatalln("fail to dial ", err)
		return nil
	}

	return conn
}
