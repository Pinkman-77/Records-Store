package handler

import (
	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
)

func (h *Handler) createArtist(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userID, err := h.services.Creator.GetUserIDByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	artist := recordsrestapi.Artist{
		Name:   input.Name,
		UserID: userID,
	}

	artistID, err := h.services.Creator.CreateArtist(artist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"artist_id": artistID})
}


func (h *Handler) getAllArtists(c *gin.Context) {
	artists, err := h.services.Creator.GetAllArtists()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"artists": artists})
}

func (h *Handler) getArtist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	artist, err := h.services.Creator.GetArtist(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Artist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"artist": artist})
}

func (h *Handler) updateArtist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
	}

	var updatedArtist recordsrestapi.Artist
	if err := c.BindJSON(&updatedArtist); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	err = h.services.Creator.UpdateArtist(id, updatedArtist)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artist updated successfully"})
}

func (h *Handler) deleteArtist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
			return
	}

	err = h.services.Creator.DeleteArtist(id)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Artist deleted successfully"})
}
