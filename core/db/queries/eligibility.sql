-- name: InsertEligibility :one
INSERT INTO eligibility (campaign_id, wallet_address, is_eligible, total_points, token_amount, merkle_proof)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEligibility :one
SELECT * FROM eligibility
WHERE campaign_id = $1 AND wallet_address = $2;

-- name: ListEligibleByCampaign :many
SELECT * FROM eligibility
WHERE campaign_id = $1 AND is_eligible = TRUE;

-- name: MarkClaimed :one
UPDATE eligibility
SET claim_status = 'claimed', claim_tx_hash = $3, claimed_at = NOW()
WHERE campaign_id = $1 AND wallet_address = $2
RETURNING *;

-- name: UpdateClaimStatus :one
UPDATE eligibility
SET claim_status = $3
WHERE campaign_id = $1 AND wallet_address = $2
RETURNING *;

-- name: ClaimRateAggregate :one
SELECT
    COUNT(*) FILTER (WHERE is_eligible = TRUE)                     AS total_eligible,
    COUNT(*) FILTER (WHERE claim_status = 'claimed')               AS total_claimed,
    COUNT(*) FILTER (WHERE claim_status = 'pending')               AS total_pending
FROM eligibility
WHERE campaign_id = $1;