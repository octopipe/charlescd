package server

import (
	"context"
	"log"
	"net"

	"github.com/go-logr/logr"
	"github.com/octopipe/charlescd/internal/butler/errs"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
	"google.golang.org/grpc"
)

type server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func serverIntercepetor(logger logr.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h, err := handler(ctx, req)
		if err != nil {
			return nil, errs.NewGRPCError(logger, err)
		}

		return h, nil
	}
}

func NewServer(
	logger logr.Logger,
	circleServer pbv1.CircleServiceServer,
	resourceServer pbv1.ResourceServiceServer,
) *server {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(serverIntercepetor(logger)),
	)
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
