-- name: CreateEvent :one
INSERT INTO events (
    id,
    timestamp,
    source,
    source_ip,
    destination_ip,
    source_port,
    destination_port,
    protocol,
    domain,
    url,
    file_hash,
    username,
    attack_type,
    severity,
    confidence,
    raw_payload
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events
WHERE id = $1;

-- name: ListEvents :many
SELECT * FROM events
ORDER BY timestamp DESC
LIMIT $1 OFFSET $2;

-- name: CountEvents :one
SELECT COUNT(*) FROM events;
