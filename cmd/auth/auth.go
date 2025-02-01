package main


import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Pinkman-77/records-restapi/pkg/config"
	"github.com/Pinkman-77/records-restapi/pkg/app"
	"github.com/Pinkman-77/records-restapi/pkg/storage/postgres"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogs(cfg.Env)

	log.Info(
		"starting the app",
		slog.Any("cfg", cfg.Env),
		slog.Int("Port", cfg.Grpc.Port),
	)

	// Initialize PostgreSQL storage
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.Name, cfg.Database.SslMode,
	)

	storage, err := postgres.New(dsn)
	if err != nil {
		log.Error("failed to connect to PostgreSQL", slog.Any("error", err))
		os.Exit(1)
	}

	application := app.NewApp(log, cfg.Grpc.Port, storage, cfg.TokenTll)

	if err := application.Run(); err != nil {
		log.Error("failed to start gRPC server", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server running", slog.Int("port", cfg.Grpc.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sign := <-quit

	log.Info("signal received", slog.Any("sign", sign))
	
	application.Stop()
}


func setupLogs(env string) *slog.Logger {

	const (
		envLocal = "local"
		envDev   = "dev"
		envProd  = "prod"
	)

	var log *slog.Logger

	switch env {
	case envLocal:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}