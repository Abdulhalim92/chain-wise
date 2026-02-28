package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chainwise/platform/config"
	"chainwise/platform/health"
	"chainwise/platform/logger"
	"chainwise/platform/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		// логгер ещё не создан — выводим в stderr и выходим
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(1)
	}
	log := logger.New(logger.Options{Level: cfg.LogLevel, Format: logger.Format(cfg.LogFormat), AddSource: cfg.LogAddSource, Service: "gateway"})

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/health", health.Handler())
	handler := middleware.RequestID(middleware.Recovery(log.Logger, middleware.Logging(log.Logger, mux)))

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: handler}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "error", err)
	}
}

const shutdownTimeout = 10 * time.Second

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "gateway")
}
