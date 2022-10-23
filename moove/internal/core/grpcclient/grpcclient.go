package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
)

type Client struct {
	CircleClient pbv1.CircleServiceClient
}

func NewGrpcClient() (Client, error) {
	conn, err := grpc.Dial(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return Client{}, err
	}

	circleClient := pbv1.NewCircleServiceClient(conn)

	grpcClient := Client{
		CircleClient: circleClient,
	}

	return grpcClient, nil
}
