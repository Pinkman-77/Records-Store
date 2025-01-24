package handler

import (
	"net/http"

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
       
}

func (h *Handler) getAllRecords(c *gin.Context) {

}

func (h *Handler) updateRecord(c *gin.Context) {

}

func (h *Handler) deleteRecord(c *gin.Context) {

}