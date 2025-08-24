package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerMetricsInterceptor returns a gRPC server interceptor that records metrics
func UnaryServerMetricsInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		code := codes.OK.String()
		if err != nil {
			st, _ := status.FromError(err)
			if st != nil {
				code = st.Code().String()
			}
		}

		GRPCRequestsTotal.WithLabelValues(serviceName, info.FullMethod, code).Inc()
		GRPCRequestDuration.WithLabelValues(serviceName, info.FullMethod).Observe(duration)

		return resp, err
	}
}
