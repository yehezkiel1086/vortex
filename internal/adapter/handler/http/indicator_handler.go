package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type IndicatorHandler struct {
	indicatorRepo        port.IndicatorRepository
	investigationService port.InvestigationService
}

func NewIndicatorHandler(indicatorRepo port.IndicatorRepository, investigationService port.InvestigationService) *IndicatorHandler {
	return &IndicatorHandler{
		indicatorRepo:        indicatorRepo,
		investigationService: investigationService,
	}
}

func (h *IndicatorHandler) ListIndicators(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	indType := c.Query("type")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var indicators []*domain.Indicator
	var err error

	if indType != "" {
		indicators, err = h.indicatorRepo.ListByType(c.Request.Context(), domain.IndicatorType(indType), int32(limit), int32(offset))
	} else {
		indicators, err = h.indicatorRepo.List(c.Request.Context(), int32(limit), int32(offset))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list indicators", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   indicators,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *IndicatorHandler) GetIndicatorDetails(c *gin.Context) {
	indType := c.Param("type")
	val := c.Param("value")

	indicator, observations, enrichments, relationships, err := h.investigationService.GetIndicatorDetails(
		c.Request.Context(),
		domain.IndicatorType(indType),
		val,
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "indicator not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get indicator details", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"indicator":     indicator,
		"observations":  observations,
		"enrichments":   enrichments,
		"relationships": relationships,
	})
}
