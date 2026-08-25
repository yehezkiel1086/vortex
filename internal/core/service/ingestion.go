package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yehezkiel1086/vortex/internal/core/domain"
	"github.com/yehezkiel1086/vortex/internal/core/port"
)

type IngestionService struct {
	eventRepo      port.EventRepository
	eventPublisher port.EventPublisher
}

func NewIngestionService(eventRepo port.EventRepository, eventPublisher port.EventPublisher) port.IngestionService {
	return &IngestionService{
		eventRepo:      eventRepo,
		eventPublisher: eventPublisher,
	}
}

func (s *IngestionService) IngestEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if event == nil {
		return nil, domain.ErrInvalidInput
	}

	// 1. Validation & Normalization
	if strings.TrimSpace(event.Source) == "" {
		return nil, fmt.Errorf("%w: source is required", domain.ErrInvalidInput)
	}

	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	// Severity normalization
	switch strings.ToLower(string(event.Severity)) {
	case string(domain.SeverityInformational), string(domain.SeverityLow),
		string(domain.SeverityMedium), string(domain.SeverityHigh), string(domain.SeverityCritical):
		event.Severity = domain.Severity(strings.ToLower(string(event.Severity)))
	default:
		event.Severity = domain.SeverityInformational
	}

	// Confidence clamping [0.0 - 1.0]
	if event.Confidence < 0.0 {
		event.Confidence = 0.0
	} else if event.Confidence > 1.0 {
		if event.Confidence <= 100.0 {
			event.Confidence = event.Confidence / 100.0 // Normalize percentage if given as 0-100
		} else {
			event.Confidence = 1.0
		}
	}

	// Clean fields
	event.SourceIP = strings.TrimSpace(event.SourceIP)
	event.DestinationIP = strings.TrimSpace(event.DestinationIP)
	event.Domain = strings.TrimSpace(strings.ToLower(event.Domain))
	event.URL = strings.TrimSpace(event.URL)
	event.FileHash = strings.TrimSpace(strings.ToLower(event.FileHash))
	event.Username = strings.TrimSpace(event.Username)
	event.AttackType = strings.TrimSpace(strings.ToLower(event.AttackType))

	// 2. Persist in database
	savedEvent, err := s.eventRepo.Create(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to persist event: %w", err)
	}

	// 3. Publish to message broker for background processing pipeline
	if s.eventPublisher != nil {
		if err := s.eventPublisher.PublishEvent(ctx, savedEvent); err != nil {
			// Log/return or propagate publisher error
			return savedEvent, fmt.Errorf("event saved but failed to publish to queue: %w", err)
		}
	}

	return savedEvent, nil
}
