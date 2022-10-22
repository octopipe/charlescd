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

	pbv1.CircleServiceServer
}

func NewServer(circleServer pbv1.CircleServiceServer) *server {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	s := &server{
		grpcServer:          grpcServer,
		listener:            lis,
		CircleServiceServer: circleServer,
	}
	pbv1.RegisterCircleServiceServer(grpcServer, s)

	return s
}

func (s server) Start() error {
	return s.grpcServer.Serve(s.listener)
}
