package app

import (
	"log/slog"
	"time"
	"github.com/Pinkman-77/records-restapi/pkg/grpc_service/auth"
	"github.com/Pinkman-77/records-restapi/pkg/app/grpc"

	"github.com/Pinkman-77/records-restapi/pkg/storage/postgres"
)


type gRPCApp struct {
	gRPCServer *grpcapp.App
}

func NewApp(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTLL time.Duration,
) *gRPCApp {
	storage, err := postgres.New(storagePath)

	if err != nil {
		panic(err)
	}

	authService := auth.NewAuth(log, storage, storage, storage, tokenTLL)

	gRPCServer := grpcapp.New(log, authService, port)

	return &gRPCApp{
		gRPCServer: gRPCServer,
	}
}

func (a *gRPCApp) Run() error {
	return a.gRPCServer.Run()
}

func (a *gRPCApp) Stop() {
	a.gRPCServer.Stop()
}