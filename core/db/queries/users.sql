-- name: CreateUser :one
INSERT INTO users (wallet_address)
VALUES ($1)
ON CONFLICT (wallet_address) DO NOTHING
RETURNING *;

-- name: GetUser :one 
SELECT * FROM users WHERE wallet_address = $1;

-- name: UpdateUser :one
UPDATE users
SET display_name = $2, bio = $3, avatar_url = $4, updated_at = NOW()
WHERE wallet_address = $1
RETURNING *;

-- name: UpsertTwitter :one
UPDATE users
SET twitter_id = $2, twitter_handle = $3, updated_at = NOW()
WHERE wallet_address = $1
RETURNING *;

-- name: UpsertDiscord :one
UPDATE users
SET discord_id = $2, discord_handle = $3, updated_at = NOW()
WHERE wallet_address = $1
RETURNING *;