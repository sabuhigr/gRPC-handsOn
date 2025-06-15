package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	newsv1 "github.com/sabuhigr/grpc-demo/api/news/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func init() {
	// Configure log package as json
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
}

func myUnaryInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	log.Printf("Calling method %s on remote server: %v", method, cc.Target())
	return invoker(ctx, method, req, reply, cc, opts...)
}

func main() {
	conn, err := grpc.NewClient(
		"127.0.0.1:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"load_balancing_config": {"pick_first":{}}}`), //grpc.WithDefaultServiceConfig(`{"load_balancing_config": {"round_robin":{}}}`)
		grpc.WithConnectParams(
			grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 1.6,
					MaxDelay:   120 * time.Second,
				},
				MinConnectTimeout: 5 * time.Second,
			},
		),
		grpc.WithUnaryInterceptor(myUnaryInterceptor),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			}),
	)

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := newsv1.NewNewsServiceClient(conn)

	log.Info("Starting to creating news")
	res, err := client.CreateNews(context.Background(), &newsv1.CreateNewsRequest{
		Id:      uuid.New().String(),
		Author:  "Sabu",
		Title:   "Test",
		Summary: "Test",
		Content: "Test",
		Source:  "https://sabuhi.grpc.github.io",
		Tags:    []string{"Test"},
	})

	if err != nil {
		log.Fatalf("failed to create news: %v", err)
	}

	log.WithFields(
		log.Fields{
			"status":   "successfully",
			"response": res,
		},
	).Infof(
		"News created successfully",
	)

	log.Infof("Starting to get news for ID: %v", res.Id)
	get_new_res, err := client.GetNews(context.Background(), &newsv1.GetNewsRequest{
		Id: res.Id,
	})

	if err != nil {
		log.Fatalf("failed to get news: %v", err)
	}

	log.WithFields(
		log.Fields{
			"status":   "successfully",
			"response": get_new_res,
		},
	).Infof(
		"News got successfully",
	)

}
