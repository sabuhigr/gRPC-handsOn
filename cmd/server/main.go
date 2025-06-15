package main

import (
	"context"
	"net"

	news1 "github.com/sabuhigr/grpc-demo/api/news/v1"
	ingrpc "github.com/sabuhigr/grpc-demo/internal/grpc"
	"github.com/sabuhigr/grpc-demo/internal/memstore"
	"github.com/sabuhigr/grpc-demo/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func init() {
	// Configure log package as json
	log.SetFormatter(&log.JSONFormatter{})
}

func unaryMetadataInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}

	// Example: Check for token
	if values := md["authorization"]; len(values) > 0 && values[0] == types.Static_token {
		log.Info("Successfully authenticated")
	} else {
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	// Continue with actual handler
	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(unaryMetadataInterceptor),
	)
	news1.RegisterNewsServiceServer(srv, ingrpc.NewServer(memstore.New()))
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthSrv)

	log.Info("Starting gRPC server on port :8080")

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
