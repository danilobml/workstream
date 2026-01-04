package grpc

import (
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateGrpcClient(grpcAddr string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, errors.New("unable to connect to gRRPC")
	}

	log.Printf("workstream-gateway - connected to gRPC at port: %s", grpcAddr)

	return conn, nil
}
