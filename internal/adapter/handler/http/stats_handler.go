package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type StatsHandler struct {
	investigationService port.InvestigationService
}

func NewStatsHandler(investigationService port.InvestigationService) *StatsHandler {
	return &StatsHandler{
		investigationService: investigationService,
	}
}

func (h *StatsHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "vortex",
	})
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.investigationService.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stats", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
