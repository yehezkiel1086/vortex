package domain

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource already exists")
	ErrInvalidInput   = errors.New("invalid input provided")
	ErrDatabase       = errors.New("database operation failed")
	ErrEnrichmentFail = errors.New("external enrichment failed")
	ErrUnauthorized   = errors.New("unauthorized access")
)
