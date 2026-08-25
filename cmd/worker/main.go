package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yehezkiel1086/vortex/internal/adapter/config"
	"github.com/yehezkiel1086/vortex/internal/adapter/enrichment/geoip"
	"github.com/yehezkiel1086/vortex/internal/adapter/enrichment/virustotal"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/rabbitmq"
	"github.com/yehezkiel1086/vortex/internal/adapter/storage/redis"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/service"
)

func main() {
	log.Println("[Worker] Starting Vortex Threat Intelligence Processing Worker...")

	// 1. Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("[Worker] Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Connect to PostgreSQL
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("[Worker] Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("[Worker] Connected to PostgreSQL.")

	// 3. Connect to Redis
	cache, err := redis.New(ctx, cfg.Cache)
	if err != nil {
		log.Printf("[Worker] Warning: Redis connection failed: %v. Running without Redis cache.", err)
	} else {
		defer cache.Close()
		log.Println("[Worker] Connected to Redis.")
	}

	// 4. Connect to RabbitMQ
	mq, err := rabbitmq.New(cfg.Rabbitmq)
	if err != nil {
		log.Fatalf("[Worker] Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()
	log.Println("[Worker] Connected to RabbitMQ.")

	// 5. Initialize Threat Intel Clients
	geoIPClient := geoip.New("", 5*time.Second)
	vtClient := virustotal.New(os.Getenv("VIRUSTOTAL_API_KEY"), "", 5*time.Second)

	// 6. Initialize Repositories
	indRepo := repository.NewIndicatorRepository(db.Pool)
	obsRepo := repository.NewObservationRepository(db.Pool)
	enrRepo := repository.NewEnrichmentRepository(db.Pool)
	relRepo := repository.NewRelationshipRepository(db.Pool)
	alertRepo := repository.NewAlertRepository(db.Pool)

	// 7. Initialize Pipeline Engine Services
	extractor := service.NewExtractorService(false)
	detector := service.NewDetectionService(obsRepo, indRepo)
	enricher := service.NewEnrichmentService(enrRepo, cache, geoIPClient, vtClient)
	correlator := service.NewCorrelationService(relRepo, indRepo)
	riskScorer := service.NewRiskScoringService(indRepo)
	alertEngine := service.NewAlertService(alertRepo)

	// 8. Start Background Consumer Loop
	consumer := rabbitmq.NewConsumer(mq)

	go func() {
		err := consumer.StartConsuming(ctx, func(c context.Context, event *domain.Event) error {
			log.Printf("[Worker] Processing event %s from source '%s' (IP: %s, Attack: %s)",
				event.ID, event.Source, event.SourceIP, event.AttackType)

			// Step 1: Extract IOCs
			indicators, err := extractor.ExtractIndicators(c, event)
			if err != nil {
				log.Printf("[Worker] Extraction error for event %s: %v", event.ID, err)
			}

			// Step 2: Run Detection Rules
			observations, err := detector.Detect(c, event)
			if err != nil {
				log.Printf("[Worker] Detection error for event %s: %v", event.ID, err)
			}

			// Step 3: Pipeline per Indicator
			for _, ind := range indicators {
				// 3a. Save/Upsert Indicator
				savedInd, err := indRepo.Upsert(c, ind)
				if err != nil {
					log.Printf("[Worker] Error saving indicator %s: %v", ind.Value, err)
					continue
				}

				// 3b. Threat Intel Enrichment (GeoIP + VirusTotal)
				enrichments, err := enricher.EnrichIndicator(c, savedInd)
				if err != nil {
					log.Printf("[Worker] Enrichment error for %s: %v", savedInd.Value, err)
				}

				// 3c. Correlation (Graph Relationships)
				relationships, err := correlator.Correlate(c, savedInd, event)
				if err != nil {
					log.Printf("[Worker] Correlation error for %s: %v", savedInd.Value, err)
				}

				// 3d. Multi-Factor Risk & Confidence Scoring
				riskScore, err := riskScorer.CalculateRisk(c, savedInd, observations, enrichments, relationships)
				if err != nil {
					log.Printf("[Worker] Risk scoring error for %s: %v", savedInd.Value, err)
					continue
				}

				// 3e. Alert Evaluation
				if riskScore != nil {
					alert, err := alertEngine.EvaluateAlert(c, savedInd, event, riskScore)
					if err != nil {
						log.Printf("[Worker] Alert evaluation error for %s: %v", savedInd.Value, err)
					} else if alert != nil {
						log.Printf("[Worker] 🚨 HIGH RISK ALERT! [%s] %s: %s (Risk Score: %.1f/100, Confidence: %.0f%%)",
							alert.Severity, alert.Title, savedInd.Value, alert.RiskScore, alert.Confidence*100)
					}
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[Worker] Consumer error: %v", err)
		}
	}()

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[Worker] Shutting down worker daemon...")
	cancel()
	time.Sleep(1 * time.Second)
	log.Println("[Worker] Worker exited cleanly.")
}
