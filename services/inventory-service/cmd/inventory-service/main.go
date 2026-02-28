package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
	"os"
	"os/signal"
	"syscall"

	"chainwise/platform/config"
	"chainwise/platform/health"
	"chainwise/platform/interceptors"
	"chainwise/platform/logger"
	"chainwise/platform/middleware"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(1)
	}
	log := logger.New(logger.Options{Level: cfg.LogLevel, Format: logger.Format(cfg.LogFormat), AddSource: cfg.LogAddSource, Service: "inventory-service"})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK); fmt.Fprint(w, "inventory-service") })
	mux.HandleFunc("/health", health.Handler())
	handler := middleware.RequestID(middleware.Recovery(log.Logger, middleware.Logging(log.Logger, mux)))
	httpSrv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors.UnaryRequestID, interceptors.UnaryRecovery(log.Logger), interceptors.UnaryLogging(log.Logger)),
		grpc.ChainStreamInterceptor(interceptors.StreamRequestID, interceptors.StreamRecovery(log.Logger), interceptors.StreamLogging(log.Logger)),
	)
	health.RegisterGRPC(grpcSrv)
	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		log.Error("grpc listen error", "error", err)
		os.Exit(1)
	}
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Error("grpc server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	grpcSrv.GracefulStop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "error", err)
	}
}
