package main

import (
	"net"

	news1 "github.com/sabuhigr/grpc-demo/api/news/v1"
	ingrpc "github.com/sabuhigr/grpc-demo/internal/grpc"
	"github.com/sabuhigr/grpc-demo/internal/memstore"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	// Configure log package as json
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	news1.RegisterNewsServiceServer(srv, ingrpc.NewServer(memstore.New()))
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthSrv)

	log.Info("Starting gRPC server on port :8080")

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
