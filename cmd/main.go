package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/Pinkman-77/records-restapi/pkg/handler"
	"github.com/Pinkman-77/records-restapi/repository"
	"github.com/Pinkman-77/records-restapi/service"
)

func main() {

	db, err := repository.Connect()
	if err != nil {
		log.Fatal("Coudn't connect to DB: ", err)
	}

	repo := repository.NewRepository(*db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	srv := new(recordsrestapi.Server)

	go func ()  {
		if err := srv.Start("8080", handler.InitRoutes()); err != nil {
			log.Fatal("Coudn't start server: ", err)
		}
	}()

	log.Print("The App started")

	quit := make(chan os.Signal, 1)
	
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("The App is shutting down...")

	if err := srv.Stop(context.Background()); err != nil {
		log.Fatal("Coudn't shutdown server: ", err)
	}

	if err := db.Close(); err != nil {
		log.Fatal("Coudn't close DB: ", err)
	}

}
