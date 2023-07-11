package server

import (
	"fmt"
	"log"
	"net"

	proto "github.com/SeaOfWisdom/sow_proto/lib-srv"

	"github.com/SeaOfWisdom/sow_library/src/config"
	srv "github.com/SeaOfWisdom/sow_library/src/service"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	proto.UnsafeLibraryServiceServer
	listener net.Listener
	server   *grpc.Server
	service  *srv.LibrarySrv
}

func NewGrpcServer(config *config.Config, service *srv.LibrarySrv) *GrpcServer {
	listener, err := net.Listen("tcp", config.GrpcAddress)
	if err != nil {
		log.Fatalf("could not listen to the address %s, err: %v", config.GrpcAddress, err)
	}
	serverOps := []grpc.ServerOption{
		grpc.StreamInterceptor(grpcprometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpcprometheus.UnaryServerInterceptor),
	}
	instance := &GrpcServer{
		service:  service,
		server:   grpc.NewServer(serverOps...),
		listener: listener,
	}
	grpcprometheus.EnableHandlingTimeHistogram()
	grpcprometheus.Register(instance.server)
	proto.RegisterLibraryServiceServer(instance.server, instance)

	return instance
}

func (gs *GrpcServer) Start() {
	go func() {
		if err := gs.server.Serve(gs.listener); err != nil {
			panic(fmt.Errorf("failed to serve gRPC: %v", err))
		}
	}()
}

func (s *GrpcServer) Stop() {
	s.server.Stop()
	_ = s.listener.Close()
}
