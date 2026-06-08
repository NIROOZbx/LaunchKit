-- name: CreateTask :one
INSERT INTO tasks (campaign_id, title, description, task_type, config, points, is_required, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListTasksByCampaign :many
SELECT * FROM tasks WHERE campaign_id = $1 ORDER BY sort_order ASC;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1 AND campaign_id = $2;

-- name: CountTasksByCampaign :one
SELECT COUNT(*) FROM tasks WHERE campaign_id = $1;