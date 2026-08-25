package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yehezkiel1086/vortex/internal/adapter/config"
	httpHandler "github.com/yehezkiel1086/vortex/internal/adapter/handler/http"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/rabbitmq"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/redis"
	"github.com/yehezkiel1086/vortex/internal/core/port"
	"github.com/yehezkiel1086/vortex/internal/core/service"
)

func main() {
	log.Println("[API] Starting Vortex Threat Intelligence API...")

	// 1. Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("[API] Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Connect to PostgreSQL
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("[API] Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("[API] Connected to PostgreSQL.")

	// 3. Connect to Redis (Optional cache layer)
	cache, err := redis.New(ctx, cfg.Cache)
	if err != nil {
		log.Printf("[API] Warning: Redis connection failed: %v. Running without Redis cache.", err)
	} else {
		defer cache.Close()
		log.Println("[API] Connected to Redis.")
	}

	// 4. Connect to RabbitMQ (Message broker)
	var eventPublisher port.EventPublisher
	mq, err := rabbitmq.New(cfg.Rabbitmq)
	if err != nil {
		log.Printf("[API] Warning: RabbitMQ connection failed: %v. Running without message queue dispatch.", err)
	} else {
		defer mq.Close()
		eventPublisher = rabbitmq.NewEventPublisher(mq)
		log.Println("[API] Connected to RabbitMQ.")
	}

	// 5. Initialize Repositories
	eventRepo := repository.NewEventRepository(db.Pool)
	indicatorRepo := repository.NewIndicatorRepository(db.Pool)
	obsRepo := repository.NewObservationRepository(db.Pool)
	enrRepo := repository.NewEnrichmentRepository(db.Pool)
	relRepo := repository.NewRelationshipRepository(db.Pool)
	alertRepo := repository.NewAlertRepository(db.Pool)

	// 6. Initialize Services
	ingestionSvc := service.NewIngestionService(eventRepo, eventPublisher)
	investigationSvc := service.NewInvestigationService(indicatorRepo, obsRepo, enrRepo, relRepo, eventRepo, alertRepo)
	alertSvc := service.NewAlertService(alertRepo)

	// 7. Setup Router
	router := httpHandler.NewRouter(httpHandler.Config{
		IngestionService:     ingestionSvc,
		InvestigationService: investigationSvc,
		AlertService:         alertSvc,
		EventRepo:            eventRepo,
		IndicatorRepo:        indicatorRepo,
		AllowedOrigins:       cfg.HTTP.AllowedOrigins,
	})

	portStr := cfg.HTTP.Port
	if portStr == "" {
		portStr = "8080"
	}
	hostStr := cfg.HTTP.Host
	if hostStr == "" {
		hostStr = "0.0.0.0"
	}
	addr := fmt.Sprintf("%s:%s", hostStr, portStr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 8. Start server in goroutine
	go func() {
		log.Printf("[API] Server listening on http://%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[API] Server error: %v", err)
		}
	}()

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[API] Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[API] Server forced to shutdown: %v", err)
	}

	log.Println("[API] Server exited cleanly.")
}
