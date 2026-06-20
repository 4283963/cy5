package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ultrasonic-backend/services"
)

type HistoryHandler struct {
	store *services.HistoryStore
}

func NewHistoryHandler(hs *services.HistoryStore) *HistoryHandler {
	return &HistoryHandler{store: hs}
}

type ListHistoryResponse struct {
	Records []services.HistorySummary `json:"records"`
}

func (h *HistoryHandler) ListHistory(c *gin.Context) {
	records, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if records == nil {
		records = []services.HistorySummary{}
	}
	c.JSON(http.StatusOK, ListHistoryResponse{Records: records})
}

func (h *HistoryHandler) GetHistory(c *gin.Context) {
	id := c.Param("id")
	record, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
