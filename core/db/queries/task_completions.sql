-- name: UpsertCompletion :one
INSERT INTO task_completions (campaign_id, task_id, wallet_address, status, proof, failure_reason, verified_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (task_id, wallet_address) DO UPDATE
    SET status = $4, proof = $5, failure_reason = $6, verified_at = $7, updated_at = NOW()
RETURNING *;

-- name: GetCompletion :one
SELECT * FROM task_completions
WHERE task_id = $1 AND wallet_address = $2;

-- name: ListCompletionsByWalletAndCampaign :many
SELECT tc.*, t.points, t.is_required FROM task_completions tc
JOIN tasks t ON t.id = tc.task_id
WHERE tc.campaign_id = $1 AND tc.wallet_address = $2;

-- name: GetCompletionSpeedSeconds :one
SELECT EXTRACT(EPOCH FROM (MAX(verified_at) - MIN(verified_at)))::INT AS seconds
FROM task_completions
WHERE campaign_id = $1 AND wallet_address = $2 AND status = 'verified';