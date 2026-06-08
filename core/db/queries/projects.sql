-- name: CreateProject :one
INSERT INTO projects (owner_address, name, description, logo_url, website_url, twitter_url, discord_url, token_contract, chain, treasury_wallet)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: ListProjectsByOwner :many
SELECT * FROM projects WHERE owner_address = $1 ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, description = $3, logo_url = $4, website_url = $5,
    twitter_url = $6, discord_url = $7, updated_at = NOW()
WHERE id = $1 AND owner_address = $8
RETURNING *;

-- name: SetProjectAPIKey :one
UPDATE projects SET api_key = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;