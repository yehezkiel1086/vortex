-- name: UpsertRelationship :one
INSERT INTO relationships (
    id,
    source_indicator_id,
    target_indicator_id,
    relationship_type,
    confidence,
    first_seen,
    last_seen
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (source_indicator_id, target_indicator_id, relationship_type) DO UPDATE
SET
    last_seen = EXCLUDED.last_seen,
    confidence = GREATEST(relationships.confidence, EXCLUDED.confidence)
RETURNING *;

-- name: GetRelationshipsByIndicatorID :many
SELECT * FROM relationships
WHERE source_indicator_id = $1 OR target_indicator_id = $1
ORDER BY last_seen DESC;
