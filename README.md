# News gRPC API

gRPC API for managing news articles, built in Go.  
This project demonstrates clean Go architecture, in-memory storage, robust validation, and advanced gRPC patterns including metadata, interceptors, and error details.

---

## Features

- **gRPC API** for creating, retrieving, and streaming news articles.
- **In-memory storage** for fast prototyping and testing.
- **Protobuf-based schema** ([proto/news/v1/news.proto](proto/news/v1/news.proto), [proto/news/v1/service.proto](proto/news/v1/service.proto)).
- **Advanced validation** and error reporting with rich gRPC error details.
- **Authentication** via gRPC metadata (token-based).
- **Comprehensive linting** and formatting with [golangci-lint](https://golangci-lint.run/).
- **Health checks** via gRPC Health API.
- **Client and server** implementations ([cmd/client/main.go](cmd/client/main.go), [cmd/server/main.go](cmd/server/main.go)).

---


---

## API Overview

### Protobuf Service

See [proto/news/v1/service.proto](proto/news/v1/service.proto):

```proto
service NewsService {
  rpc CreateNews(CreateNewsRequest) returns (CreateNewsResponse);
  rpc GetNews(GetNewsRequest) returns (GetNewsResponse);
  rpc GetAll(google.protobuf.Empty) returns (stream GetNewsResponse);
}
```

### Install Tools
```
make install-tools
```

### Generate Protobuf Code
```
make generate-proto
```

### Lint 
```
make lint
```



### Authentication
All requests have authentication.
- On server-side there is interceptor for validation token.
- On client-side, authorization is added to context that it passes to server.


### Error Handling
Uses gRPC status codes and rich error details (see internal/grpc/server.go).
Validation errors return INVALID_ARGUMENT with detailed descriptions.

### Development & Linting
Lint and format with make lint and make format.
Protobuf linting and breaking change detection via Buf (make lint-proto, make lint-breaking).


### Extending
Add new fields to proto/news/v1/news.proto and regenerate code.
Implement persistent storage by extending memstore.Store.
