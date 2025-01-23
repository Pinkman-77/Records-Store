package handler

import (
        "github.com/gin-gonic/gin"
        "net/http"

        "github.com/Pinkman-77/records-restapi"
)

func (h *Handler) createRecord(c *gin.Context) {
        var input recordsrestapi.Record
        if err := c.BindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        newRecord, err := h.services.Record.CreateRecord(input)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusCreated, newRecord)
}

func (h *Handler) getRecord(c *gin.Context) {
       
}

func (h *Handler) getAllRecords(c *gin.Context) {

}

func (h *Handler) updateRecord(c *gin.Context) {

}

func (h *Handler) deleteRecord(c *gin.Context) {

}