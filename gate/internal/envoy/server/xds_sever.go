package server

import (
	"context"
	"fmt"
	"log"
	"net"

	cluster "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	grpcMaxConcurrentStreams = 1000000
)

type XdsServer struct {
	managementPort uint
	ctx            context.Context
	server         xds.Server
	snapshotCache  cache.SnapshotCache
}

type healthServer struct {
	health.UnimplementedHealthServer
}

func NewXdsServer(managementPort uint, callbacks xds.Callbacks) XdsServer {
	ctx := context.Background()
	snapshotCache := cache.NewSnapshotCache(true, cache.IDHash{}, nil)
	srv := xds.NewServer(ctx, snapshotCache, callbacks)

	return XdsServer{
		ctx:           ctx,
		server:        srv,
		snapshotCache: snapshotCache,
	}

}

func (x XdsServer) Start() error {
	port := x.managementPort
	server := x.server

	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	health.RegisterHealthServer(grpcServer, healthServer{})
	cluster.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	listener.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	route.RegisterRouteDiscoveryServiceServer(grpcServer, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	return grpcServer.Serve(lis)

}
