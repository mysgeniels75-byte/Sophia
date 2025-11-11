// Package main implements the OSE Advisory Service API Gateway.
//
// This server orchestrates all backend components (Pattern Graph, Template Engine)
// to provide a unified advisory interface via gRPC and HTTP/REST.
//
// Architecture:
//   - gRPC on port 50051 for CLI communication (primary interface)
//   - HTTP on port 8080 for health checks and metrics
//   - Structured logging via zap
//   - Prometheus metrics
//   - Graceful shutdown on SIGTERM/SIGINT
//
// Configuration:
//   - Environment variables (see internal/config/config.go)
//   - Defaults work for local development
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mysgeniels75-byte/ose-api/internal/handlers"
	pb "github.com/mysgeniels75-byte/ose-api/api/proto/advisory/v1"
)

// Version information (set via ldflags during build)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Initialize structured logger
	logger, err := initLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting OSE Advisory Service",
		zap.String("version", Version),
		zap.String("commit", GitCommit),
		zap.String("build_time", BuildTime),
	)

	// Load configuration
	// Note: Actual config package import will be added when we have the full structure
	cfg := &Config{
		GRPCAddress: getEnv("OSE_GRPC_ADDRESS", ":50051"),
		HTTPAddress: getEnv("OSE_HTTP_ADDRESS", ":8080"),
		LogLevel:    getEnv("OSE_LOG_LEVEL", "info"),
	}

	logger.Info("Configuration loaded",
		zap.String("grpc_address", cfg.GRPCAddress),
		zap.String("http_address", cfg.HTTPAddress),
		zap.String("log_level", cfg.LogLevel),
	)

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(logger),
			metricsInterceptor(),
			recoveryInterceptor(logger),
		),
	)

	// Register Advisory Service handler
	advisoryHandler := handlers.NewAdvisoryHandler(logger)
	pb.RegisterAdvisoryServiceServer(grpcServer, advisoryHandler)

	// Enable gRPC reflection for debugging (grpcurl, grpc_cli)
	reflection.Register(grpcServer)

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		logger.Fatal("Failed to listen on gRPC address", zap.Error(err))
	}

	grpcErrChan := make(chan error, 1)
	go func() {
		logger.Info("Starting gRPC server", zap.String("address", cfg.GRPCAddress))
		grpcErrChan <- grpcServer.Serve(grpcListener)
	}()

	// Start HTTP server (health checks + metrics)
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/health", healthCheckHandler)
	httpMux.HandleFunc("/ready", readinessHandler)
	httpMux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      httpMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	httpErrChan := make(chan error, 1)
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", cfg.HTTPAddress))
		httpErrChan <- httpServer.ListenAndServe()
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-grpcErrChan:
		logger.Fatal("gRPC server failed", zap.Error(err))
	case err := <-httpErrChan:
		if err != http.ErrServerClosed {
			logger.Fatal("HTTP server failed", zap.Error(err))
		}
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	}

	// Graceful shutdown
	logger.Info("Shutting down gracefully...")

	// Stop accepting new gRPC requests (waits for in-flight to complete)
	grpcServer.GracefulStop()

	// Shutdown HTTP server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	logger.Info("Shutdown complete")
}

// Config holds server configuration
type Config struct {
	GRPCAddress string
	HTTPAddress string
	LogLevel    string
}

// initLogger creates a production-ready structured logger
func initLogger() (*zap.Logger, error) {
	logLevel := getEnv("OSE_LOG_LEVEL", "info")

	var cfg zap.Config
	if logLevel == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	return cfg.Build()
}

// healthCheckHandler returns 200 OK if server is running (liveness probe)
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// readinessHandler returns 200 OK if server is ready to accept traffic
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check backend dependencies (Pattern Graph, etc.)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}

// getEnv retrieves environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ═══════════════════════════════════════════════════════════════════════
// MIDDLEWARE INTERCEPTORS (inline for Week 3)
// ═══════════════════════════════════════════════════════════════════════

func loggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.Error("RPC failed", fields...)
		} else {
			logger.Info("RPC succeeded", fields...)
		}

		return resp, err
	}
}

func metricsInterceptor() grpc.UnaryServerInterceptor {
	// TODO: Implement Prometheus metrics in Week 3 continuation
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		return handler(ctx, req)
	}
}

func recoveryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Panic recovered in RPC handler",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
				)
				err = fmt.Errorf("internal server error: %v", r)
			}
		}()

		return handler(ctx, req)
	}
}
