package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type Config struct {
	IngestionService     port.IngestionService
	InvestigationService port.InvestigationService
	AlertService         port.AlertService
	EventRepo            port.EventRepository
	IndicatorRepo        port.IndicatorRepository
	AllowedOrigins       string
}

func NewRouter(cfg Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Simple CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Handlers
	eventHandler := NewEventHandler(cfg.IngestionService, cfg.EventRepo)
	indicatorHandler := NewIndicatorHandler(cfg.IndicatorRepo, cfg.InvestigationService)
	alertHandler := NewAlertHandler(cfg.AlertService)
	statsHandler := NewStatsHandler(cfg.InvestigationService)

	// Liveness probe
	r.GET("/health", statsHandler.Health)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Events
		v1.POST("/events", eventHandler.IngestEvent)
		v1.GET("/events", eventHandler.ListEvents)
		v1.GET("/events/:id", eventHandler.GetEventByID)

		// Indicators
		v1.GET("/indicators", indicatorHandler.ListIndicators)
		v1.GET("/indicators/:type/:value", indicatorHandler.GetIndicatorDetails)

		// Alerts
		v1.GET("/alerts", alertHandler.ListAlerts)
		v1.GET("/alerts/:id", alertHandler.GetAlertByID)
		v1.PATCH("/alerts/:id", alertHandler.UpdateAlertStatus)

		// Statistics
		v1.GET("/stats", statsHandler.GetStats)
	}

	return r
}
