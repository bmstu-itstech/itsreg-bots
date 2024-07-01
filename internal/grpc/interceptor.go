package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

// LoggingInterceptor is a gRPC server interceptor that logs incoming requests using Zap sugared logger.
func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		log = log.With(
			slog.With("method", info.FullMethod),
			slog.With("request", req))
		log.Info("incoming request")

		resp, err := handler(ctx, req)
		duration := time.Since(startTime)

		log = log.With(
			slog.With("method", info.FullMethod),
			slog.With("duration", duration))

		if err != nil {
			log.Error("error processing request", "err", err)
		} else {
			log.Info("processed request")
		}

		return resp, err
	}
}
