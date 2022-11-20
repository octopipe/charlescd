package server

import (
	"log"
	"net"

	"github.com/go-logr/logr"
	pbv1 "github.com/octopipe/charlescd/compass/pb/v1"
	"google.golang.org/grpc"
)

type server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(
	logger logr.Logger,
	resourceServer pbv1.ResourceServiceServer,
) *server {
	lis, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	s := &server{
		grpcServer: grpcServer,
		listener:   lis,
	}
	pbv1.RegisterResourceServiceServer(grpcServer, resourceServer)

	return s
}

func (s server) Start() error {
	return s.grpcServer.Serve(s.listener)
}
