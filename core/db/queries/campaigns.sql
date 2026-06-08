-- name: CreateCampaign :one
INSERT INTO campaigns (project_id, name, description, banner_url, chain, token_contract, total_allocation, reward_type, reward_config, eligibility_rules, vesting_type, vesting_days, claim_window_days, gas_sponsored, starts_at, ends_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: GetCampaign :one
SELECT * FROM campaigns WHERE id = $1;

-- name: ListActiveCampaigns :many
SELECT * FROM campaigns
WHERE status = 'active'
  AND ($1::VARCHAR = '' OR chain = $1)
  AND ($2::VARCHAR = '' OR reward_type::TEXT = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateCampaign :one
UPDATE campaigns
SET name = $2, description = $3, banner_url = $4, reward_config = $5,
    eligibility_rules = $6, starts_at = $7, ends_at = $8, updated_at = NOW()
WHERE id = $1 AND status = 'draft'
RETURNING *;

-- name: UpdateCampaignStatus :one
UPDATE campaigns SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SetCampaignMerkleData :one
UPDATE campaigns
SET merkle_root = $2, contract_address = $3, status = 'distributing', updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ExtendCampaignDeadline :one
UPDATE campaigns SET ends_at = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;