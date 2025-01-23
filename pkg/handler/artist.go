package handler

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createArtist(c *gin.Context) {
	var input recordsrestapi.Artist

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	artist, err := h.services.Creator.CreateArtist(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]interface{}{
		"artist": artist,
	})

}

func (h *Handler) getAllArtists(c *gin.Context) {
	artists, err := h.services.Creator.GetAllArtists()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, artists)
}

func (h *Handler) getArtist(c *gin.Context) {

}

func (h *Handler) updateArtist(c *gin.Context) {

}

func (h *Handler) deleteArtist(c *gin.Context) {

}