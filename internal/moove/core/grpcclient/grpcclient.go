package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbv1 "github.com/octopipe/charlescd/pb/v1"
)

type Client struct {
	ResourceClient pbv1.ResourceServiceClient
}

func NewGrpcClient() (Client, error) {
	conn, err := grpc.Dial(
		"localhost:3000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return Client{}, err
	}

	resourceClient := pbv1.NewResourceServiceClient(conn)

	grpcClient := Client{
		ResourceClient: resourceClient,
	}

	return grpcClient, nil
}
