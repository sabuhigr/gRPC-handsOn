package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	newsv1 "github.com/sabuhigr/grpc-demo/api/news/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
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

	log.Info("Starting to bulk creating news")
	for i := range 5 {
		_, err := client.CreateNews(context.Background(), &newsv1.CreateNewsRequest{
			Id:      uuid.New().String(),
			Author:  fmt.Sprintf("Test Author %v", i),
			Title:   "Test",
			Summary: "Test",
			Content: "Test",
			Source:  "https://sabuhi.grpc.github.io",
			Tags:    []string{"Test"},
		})

		if err != nil {
			log.Errorf("failed to create news: %v", err)
			return
		}
	}

	log.Infof(
		"Bulk news created successfully",
	)

	log.Info("Starting to get all news")

	allnewsRes, err := client.GetAll(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Errorf("failed to get all news: %v", err)
		return
	}

	allNews := make([]*newsv1.GetNewsResponse, 0)

	for {

		resp, err := allnewsRes.Recv()
		if err != nil {
			if err == io.EOF {
				log.WithField("news", allNews).Infof(
					"Got all news successfully",
				)
				break
			} else {
				log.Errorf("failed to get all news: %v", err)
				return
			}
		}

		allNews = append(allNews, resp)
		log.Infof("Got author [%v] news successfully.ID - %v", resp.Author, resp.Id)
	}
}
