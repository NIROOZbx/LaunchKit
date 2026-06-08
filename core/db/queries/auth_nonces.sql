-- name: UpsertNonce :one
INSERT INTO auth_nonces (wallet_address, nonce, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (wallet_address) DO UPDATE
    SET nonce = $2, expires_at = $3
RETURNING *;

-- name: GetNonce :one
SELECT * FROM auth_nonces
WHERE wallet_address = $1 AND expires_at > NOW();

-- name: DeleteNonce :exec
DELETE FROM auth_nonces WHERE wallet_address = $1;

-- name: DeleteExpiredNonces :exec
DELETE FROM auth_nonces WHERE expires_at <= NOW();