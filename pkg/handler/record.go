package handler

import (
	"net/http"
	"strconv"

	recordsrestapi "github.com/Pinkman-77/records-restapi"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createRecord(c *gin.Context) {
	var input recordsrestapi.Record

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordID, err := h.services.Record.CreateRecord(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record_id": recordID})
}

func (h *Handler) getRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	record, err := h.services.Record.GetRecord(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": record})
}

func (h *Handler) getAllRecords(c *gin.Context) {
	records, err := h.services.Record.GetAllRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}


func (h *Handler) updateRecord(c *gin.Context) {

}

func (h *Handler) deleteRecord(c *gin.Context) {

}