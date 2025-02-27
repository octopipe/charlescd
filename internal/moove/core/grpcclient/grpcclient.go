package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbv1 "github.com/octopipe/charlescd/pb/v1"
)

type Client struct {
	CircleClient   pbv1.CircleServiceClient
	ResourceClient pbv1.ResourceServiceClient
}

func NewGrpcClient() (Client, error) {
	conn, err := grpc.Dial(
		"localhost:8002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return Client{}, err
	}

	circleClient := pbv1.NewCircleServiceClient(conn)
	resourceClient := pbv1.NewResourceServiceClient(conn)

	grpcClient := Client{
		CircleClient:   circleClient,
		ResourceClient: resourceClient,
	}

	return grpcClient, nil
}
