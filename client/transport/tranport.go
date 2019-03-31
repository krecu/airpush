package transport

import "context"

// interface for GRPC/HTTP connection
type Transport interface {
	Do(ctx context.Context) ([]byte, error)
}