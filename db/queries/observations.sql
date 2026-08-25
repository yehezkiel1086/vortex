-- name: CreateObservation :one
INSERT INTO observations (
    id,
    indicator_id,
    event_id,
    attack_type,
    technique_id,
    timestamp,
    severity,
    confidence,
    source
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: ListObservationsByIndicatorID :many
SELECT * FROM observations
WHERE indicator_id = $1
ORDER BY timestamp DESC;

-- name: ListObservationsByEventID :many
SELECT * FROM observations
WHERE event_id = $1
ORDER BY timestamp DESC;
