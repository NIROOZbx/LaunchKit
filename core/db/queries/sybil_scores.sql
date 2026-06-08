-- name: InsertSybilScore :one
INSERT INTO sybil_scores (campaign_id, wallet_address, wallet_age_score, tx_count_score, speed_score, passport_score, funding_source_score, ip_cluster_score, total_score, flags, is_eligible)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetSybilScore :one
SELECT * FROM sybil_scores
WHERE campaign_id = $1 AND wallet_address = $2;

-- name: ListFlaggedByCampaign :many
SELECT * FROM sybil_scores
WHERE campaign_id = $1 AND is_eligible = FALSE
ORDER BY total_score DESC;

-- name: ScoreDistribution :many
SELECT
    CASE
        WHEN total_score < 20 THEN '0-19'
        WHEN total_score < 40 THEN '20-39'
        WHEN total_score < 60 THEN '40-59'
        WHEN total_score < 80 THEN '60-79'
        ELSE '80-100'
    END AS bucket,
    COUNT(*) AS count
FROM sybil_scores
WHERE campaign_id = $1
GROUP BY bucket ORDER BY bucket;