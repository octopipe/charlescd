package server

import (
	"log"
	"net"

	pbv1 "github.com/octopipe/charlescd/butler/pb/v1"
	"google.golang.org/grpc"
)

type server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(
	circleServer pbv1.CircleServiceServer,
	resourceServer pbv1.ResourceServiceServer,
) *server {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	s := &server{
		grpcServer: grpcServer,
		listener:   lis,
	}
	pbv1.RegisterCircleServiceServer(grpcServer, circleServer)
	pbv1.RegisterResourceServiceServer(grpcServer, resourceServer)

	return s
}

func (s server) Start() error {
	return s.grpcServer.Serve(s.listener)
}
