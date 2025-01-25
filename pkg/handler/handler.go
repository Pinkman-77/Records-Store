package handler

import (
	"github.com/Pinkman-77/records-restapi/pkg/service"
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
		artist := api.Group("/artists")
		{
			artist.GET("/:id", h.getArtist)
			artist.GET("/", h.getAllArtists)
			artist.POST("/", h.createArtist)
			artist.PUT("/:id", h.updateArtist)
			artist.DELETE("/:id", h.deleteArtist) 
		}
		record := api.Group("/records")
		{
			record.GET("/:id", h.getRecord)
			record.GET("/", h.getAllRecords)
			record.POST("/", h.createRecord)
			record.PUT("/:id", h.updateRecord)
			record.PATCH("/:id", h.patchRecord)
			record.DELETE("/:id", h.deleteRecord) 
		}
	}
	return router
}
