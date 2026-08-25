package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type EventHandler struct {
	ingestionService port.IngestionService
	eventRepo        port.EventRepository
}

func NewEventHandler(ingestionService port.IngestionService, eventRepo port.EventRepository) *EventHandler {
	return &EventHandler{
		ingestionService: ingestionService,
		eventRepo:        eventRepo,
	}
}

func (h *EventHandler) IngestEvent(c *gin.Context) {
	var event domain.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	ingested, err := h.ingestionService.IngestEvent(c.Request.Context(), &event)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ingest event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "event ingested successfully",
		"data":    ingested,
	})
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	events, err := h.eventRepo.List(c.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   events,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *EventHandler) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	event, err := h.eventRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": event})
}
