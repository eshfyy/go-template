package interceptor

import (
	"context"
	"errors"
	"go-template/internal/domain"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorUnaryInterceptor maps domain errors to gRPC status codes.
func ErrorUnaryInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}
		return nil, mapError(err, log, info.FullMethod)
	}
}

func mapError(err error, log *zap.Logger, method string) error {
	var validationErr *domain.ValidationError
	switch {
	case errors.As(err, &validationErr):
		return status.Error(codes.InvalidArgument, validationErr.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "already exists")
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid input")
	default:
		log.Error("internal error", zap.Error(err), zap.String("method", method))
		return status.Error(codes.Internal, "internal error")
	}
}
