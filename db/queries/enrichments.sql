-- name: UpsertEnrichment :one
INSERT INTO enrichments (
    id,
    indicator_id,
    provider,
    data,
    fetched_at,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (indicator_id, provider) DO UPDATE
SET
    data = EXCLUDED.data,
    fetched_at = EXCLUDED.fetched_at,
    expires_at = EXCLUDED.expires_at
RETURNING *;

-- name: GetEnrichmentsByIndicatorID :many
SELECT * FROM enrichments
WHERE indicator_id = $1
ORDER BY fetched_at DESC;

-- name: GetEnrichmentByProvider :one
SELECT * FROM enrichments
WHERE indicator_id = $1 AND provider = $2;
