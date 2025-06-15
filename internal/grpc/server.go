package grpc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	newsv1 "github.com/sabuhigr/grpc-demo/api/news/v1"
	"github.com/sabuhigr/grpc-demo/internal/memstore"
	"github.com/sabuhigr/grpc-demo/types"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type NewsStorer interface {
	Create(news *memstore.News) *memstore.News
	Get(id uuid.UUID) *memstore.News
	GetAll() []*memstore.News
}

// Server gRPC server.
type Server struct {
	newsv1.UnimplementedNewsServiceServer
	store NewsStorer
}

// NewServer creates a new gRPC server as pointer.
func NewServer(store NewsStorer) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) ErrorWithDetails(code codes.Code, errDetails types.ErrDetails) error {
	st := status.Newf(code, fmt.Sprintf("something went wrong: %v", errDetails.Message))
	v := &errdetails.PreconditionFailure_Violation{ //errDetails
		Type:        errDetails.Type,
		Subject:     errDetails.Message,
		Description: errDetails.Description,
	}
	br := &errdetails.PreconditionFailure{}
	br.Violations = append(br.Violations, v)
	st, _ = st.WithDetails(br)
	return st.Err()
}

func (s *Server) CreateNews(context context.Context, in *newsv1.CreateNewsRequest) (*newsv1.CreateNewsResponse, error) {
	log := log.WithFields(
		log.Fields{
			"request_data": in,
			"endpoint":     "GetNews",
		})

	log.Debugf("Received request from client")
	parsedNews, err := parseAndValidate(in)
	if err != nil {
		return nil, s.ErrorWithDetails(codes.InvalidArgument, types.ErrDetails{Code: 400, Message: err.Error(), Type: "invalid_argument", Description: "invalid argument"})
	} else {
		createdNews := s.store.Create(parsedNews)
		log.WithFields(
			logrus.Fields{
				"status": "successfully",
				"news":   createdNews,
			},
		).Infof("News created successfully!")
		return toNewsResponse(createdNews), nil
	}
}
func (s *Server) GetNews(context context.Context, in *newsv1.GetNewsRequest) (*newsv1.GetNewsResponse, error) {
	log := log.WithFields(
		log.Fields{
			"request_data": in,
			"endpoint":     "GetNews",
		})

	log.Debugf("Received request from client")
	parseUUID, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, s.ErrorWithDetails(codes.InvalidArgument, types.ErrDetails{Code: 400, Message: err.Error(), Type: "invalid_UUID", Description: "invalid UUID"})
	}

	log.Debugf("uuid: %v", parseUUID)

	news := s.store.Get(parseUUID)
	log.Debugf("news: %v", news)
	if news == nil {
		return nil, s.ErrorWithDetails(codes.NotFound, types.ErrDetails{Code: 500, Message: err.Error(), Type: "not_found", Description: "Not found"})
	}

	log.WithFields(
		logrus.Fields{
			"status": "successfully",
			"news":   news,
		},
	).Infof("News got from memstore successfully!")
	return toGetNewsResponse(news), nil
}

func (s *Server) GetAll(in *emptypb.Empty, stream newsv1.NewsService_GetAllServer) error {
	log := log.WithFields(
		log.Fields{
			"request_data": in,
			"endpoint":     "GetAll",
		},
	)

	log.Debugf("Received request from client")
	newsList := s.store.GetAll()
	for _, news := range newsList {
		if err := stream.Send(toGetNewsResponse(news)); err != nil {
			return err
		}
	}
	return nil
}

func parseAndValidate(in *newsv1.CreateNewsRequest) (n *memstore.News, errs error) {
	errs = nil
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if in.Author == "" {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "author cannot be empty"))
	}

	if in.Title == "" {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "title cannot be empty"))
	}

	if in.Summary == "" {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "summary cannot be empty"))
	}

	if in.Content == "" {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "content cannot be empty"))
	}

	if in.Source == "" {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "source cannot be empty"))
	}

	if len(in.Tags) == 0 {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "tags cannot be empty"))
	}

	parsedID, err := uuid.Parse(in.Id)
	if err != nil {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "id cannot be parsed"))
	}

	parsedURL, err := url.Parse(in.Source)
	if err != nil {
		errs = errors.Join(errs, status.Errorf(codes.InvalidArgument, "source cannot be parsed"))
	}

	if errs != nil {
		return nil, errs
	}

	return &memstore.News{
		ID:        parsedID,
		Author:    in.Author,
		Title:     in.Title,
		Summary:   in.Summary,
		Content:   in.Content,
		Source:    parsedURL,
		Tags:      in.Tags,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: time.Now().UTC(),
	}, nil
}

func toNewsResponse(news *memstore.News) *newsv1.CreateNewsResponse {
	if news == nil {
		return nil
	}

	return &newsv1.CreateNewsResponse{
		Id:        news.ID.String(),
		Author:    news.Author,
		Title:     news.Title,
		Summary:   news.Summary,
		Content:   news.Content,
		Source:    news.Source.String(),
		Tags:      news.Tags,
		CreatedAt: timestamppb.New(news.CreatedAt.UTC()),
		UpdatedAt: timestamppb.New(news.UpdatedAt.UTC()),
	}
}

func toGetNewsResponse(news *memstore.News) *newsv1.GetNewsResponse {
	if news == nil {
		return nil
	}

	return &newsv1.GetNewsResponse{
		Id:        news.ID.String(),
		Author:    news.Author,
		Title:     news.Title,
		Summary:   news.Summary,
		Content:   news.Content,
		Source:    news.Source.String(),
		Tags:      news.Tags,
		CreatedAt: timestamppb.New(news.CreatedAt.UTC()),
		UpdatedAt: timestamppb.New(news.UpdatedAt.UTC()),
	}
}
