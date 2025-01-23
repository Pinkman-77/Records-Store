package handler

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/gin-gonic/gin"
	"strconv"
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
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid artist ID"})
		return
	}

	artist, err := h.services.Creator.GetArtist(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Artist not found"})
		return
	}

	c.JSON(200, artist)
}

func (h *Handler) updateArtist(c *gin.Context) {

}

func (h *Handler) deleteArtist(c *gin.Context) {

}