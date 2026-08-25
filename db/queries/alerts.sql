-- name: CreateAlert :one
INSERT INTO alerts (
    id,
    indicator_id,
    event_id,
    severity,
    risk_score,
    confidence,
    title,
    description,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetAlertByID :one
SELECT * FROM alerts
WHERE id = $1;

-- name: ListAlerts :many
SELECT * FROM alerts
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListAlertsByStatus :many
SELECT * FROM alerts
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAlertStatus :one
UPDATE alerts
SET
    status = $2,
    resolved_at = CASE WHEN $2 IN ('resolved', 'false_positive') THEN now() ELSE NULL END
WHERE id = $1
RETURNING *;

-- name: CountAlertsByStatus :many
SELECT status, COUNT(*) as count
FROM alerts
GROUP BY status;
