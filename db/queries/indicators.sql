-- name: UpsertIndicator :one
INSERT INTO indicators (
    id,
    type,
    value,
    first_seen,
    last_seen,
    risk_score,
    confidence,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (type, value) DO UPDATE
SET
    last_seen = EXCLUDED.last_seen,
    risk_score = GREATEST(indicators.risk_score, EXCLUDED.risk_score),
    confidence = GREATEST(indicators.confidence, EXCLUDED.confidence),
    updated_at = now()
RETURNING *;

-- name: GetIndicatorByID :one
SELECT * FROM indicators
WHERE id = $1;

-- name: GetIndicatorByTypeValue :one
SELECT * FROM indicators
WHERE type = $1 AND value = $2;

-- name: ListIndicators :many
SELECT * FROM indicators
ORDER BY risk_score DESC, last_seen DESC
LIMIT $1 OFFSET $2;

-- name: ListIndicatorsByType :many
SELECT * FROM indicators
WHERE type = $1
ORDER BY risk_score DESC, last_seen DESC
LIMIT $2 OFFSET $3;

-- name: UpdateIndicatorRiskScore :one
UPDATE indicators
SET
    risk_score = $2,
    confidence = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CountIndicators :one
SELECT COUNT(*) FROM indicators;
