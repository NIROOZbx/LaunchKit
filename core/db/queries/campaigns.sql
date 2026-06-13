INSERT INTO campaigns (
    project_id, title, slug, description, banner_url,
    chain, token_contract, token_symbol, token_decimals,
    total_allocation, reward_type, reward_config,
    eligibility_rules, start_date, end_date, claim_window_days
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12,
    $13, $14, $15, $16
)
RETURNING *;

-- name: GetCampaignByID :one
SELECT * FROM campaigns
WHERE id = $1
AND deleted_at IS NULL;

-- name: GetCampaignBySlug :one
SELECT * FROM campaigns
WHERE slug = $1
AND deleted_at IS NULL;

-- name: ListCampaignsByProject :many
SELECT * FROM campaigns
WHERE project_id = $1
AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListActiveCampaigns :many
SELECT * FROM campaigns
WHERE status = 'active'
AND chain = COALESCE(sqlc.narg(chain), chain)
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCampaignStatus :one
UPDATE campaigns
SET status = $2
WHERE id = $1
AND deleted_at IS NULL
RETURNING *;

-- name: UpdateCampaign :one
UPDATE campaigns SET
    title             = COALESCE($2, title),
    description       = COALESCE($3, description),
    banner_url        = COALESCE($4, banner_url),
    reward_config     = COALESCE($5, reward_config),
    eligibility_rules = COALESCE($6, eligibility_rules),
    start_date        = COALESCE($7, start_date),
    end_date          = COALESCE($8, end_date),
    claim_window_days = COALESCE($9, claim_window_days)
WHERE id = $1
AND status = 'draft'
AND deleted_at IS NULL
RETURNING *;

-- name: FinalizeCampaign :one
UPDATE campaigns SET
    status         = 'distributing',
    merkle_root    = $2,
    claim_contract = $3,
    deploy_tx_hash = $4,
    finalized_at   = NOW()
WHERE id = $1
AND status = 'ended'
RETURNING *;

-- name: IncrementClaimedAmount :one
UPDATE campaigns SET
    claimed_amount = claimed_amount + $2
WHERE id = $1
RETURNING *;

-- name: SoftDeleteCampaign :exec
UPDATE campaigns
SET deleted_at = NOW()
WHERE id = $1
AND deleted_at IS NULL;