package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/bridges/otelslog"

	"goManip/jobs"
	"goManip/server"
	"goManip/store"
	"goManip/worker"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		slog.Warn("Invalid integer in env var, using default", "env", key, "value", raw, "error", err)
		return defaultValue
	}

	return value
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		slog.Warn("Invalid duration in env var, using default", "env", key, "value", raw, "error", err)
		return defaultValue
	}

	return value
}

func main() {
	numWorkers := flag.Int("numWorkers", runtime.NumCPU(), "Number of workers")
	port := flag.String("port", "8080", "port to listen on")
	address := flag.String("address", "localhost", "address to bind to")
	flag.Parse()

	jobReqs := make(chan *jobs.JobRequest)
	jobDispatcher := jobs.NewJobDispatcher(jobReqs)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	redisAddr := getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := getEnvIntOrDefault("REDIS_DB", 0)
	redisRetention := getEnvDurationOrDefault("REDIS_RETENTION_DURATION", 24*time.Hour)

	exporterEndpoint := getEnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318")
	traceShutdownTimeout := getEnvDurationOrDefault("OTEL_TRACES_SHUTDOWN_TIMEOUT", 5*time.Second)

	shutdownTracer, err := server.InitTracer(exporterEndpoint)
	if err != nil {
		slog.Error("Failed to initialize tracing", "error", err)
		os.Exit(1)
	}

	shutdownLogger, err := server.InitLogger(exporterEndpoint)
	if err != nil {
		slog.Error("Failed to initialize OTel logger", "error", err)
		os.Exit(1)
	}

	imageStore := store.NewRedisStore(redisAddr, redisUsername, redisPassword, redisDB, redisRetention)

	go func() {
		<-c
		server.GraceFullShutdown(jobDispatcher, wg, cancel)
		if err := imageStore.Close(); err != nil {
			slog.Error("Failed to close image store during shutdown", "error", err)
		}
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), traceShutdownTimeout)
		defer shutdownCancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			slog.Error("Failed to flush tracing data during shutdown", "error", err)
		}
		if err := shutdownLogger(shutdownCtx); err != nil {
			slog.Error("Failed to flush logging data during shutdown", "error", err)
		}
		os.Exit(0)
	}()

	logger := otelslog.NewLogger("gomanip")
	slog.SetDefault(logger)

	for workerId := range *numWorkers {
		slog.Info(fmt.Sprintf("Starting worker #%d", workerId+1))
		wg.Add(1)
		worker := worker.NewWorker(ctx, workerId+1, jobReqs, wg, imageStore, logger)
		go worker.Work()
	}
	e := echo.New()

	server.InitRouting(e, jobDispatcher, imageStore)
	server.Start(e, *address, *port)
}
