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

type AlertHandler struct {
	alertService port.AlertService
}

func NewAlertHandler(alertService port.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

func (h *AlertHandler) ListAlerts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	statusStr := c.Query("status")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var status *domain.AlertStatus
	if statusStr != "" {
		st := domain.AlertStatus(statusStr)
		status = &st
	}

	alerts, err := h.alertService.ListAlerts(c.Request.Context(), status, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list alerts", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   alerts,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *AlertHandler) GetAlertByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	alert, err := h.alertService.GetAlert(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get alert", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": alert})
}

type updateAlertStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *AlertHandler) UpdateAlertStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	var req updateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	updated, err := h.alertService.UpdateAlertStatus(c.Request.Context(), id, domain.AlertStatus(req.Status))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update alert status", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "alert status updated",
		"data":    updated,
	})
}
