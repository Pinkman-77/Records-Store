package main

import (
	"log"

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

	if err := srv.Start("8080", handler.InitRoutes()); err != nil {
		log.Fatal("Coudn't start server: ", err)
	}
	
}
