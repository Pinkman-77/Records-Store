package handler

import (
	"github.com/Pinkman-77/records-restapi/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		record := api.Group("/records")
		{
			record.GET("/:id", h.getRecord)
			record.GET("/", h.getAllRecords)
			record.POST("/", h.createRecord)
			record.PUT("/:id", h.updateRecord)
			record.DELETE("/:id", h.deleteRecord) 
		}
	}
	return router
}
