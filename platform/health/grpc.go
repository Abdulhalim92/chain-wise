package health

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

// RegisterGRPC регистрирует gRPC health-check на сервере.
func RegisterGRPC(srv *grpc.Server) {
	hs := health.NewServer()
	hs.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
	healthgrpc.RegisterHealthServer(srv, hs)
}
