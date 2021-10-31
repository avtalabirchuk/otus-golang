package logmiddleware

import (
	"context"
	"time"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func ApplyGRPC() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)
		code := grpc_logging.DefaultErrorToCode(err)
		latency := time.Since(startTime)

		log.Info().Msgf("[GRPC] %s %s %s", info.FullMethod, code, latency)
		return resp, err
	}
}
