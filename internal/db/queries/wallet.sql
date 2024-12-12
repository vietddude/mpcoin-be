-- name: CreateWallet :one
INSERT INTO wallets (
    user_id,
    address,
    encrypted_private_key,
    name,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetWalletsByUserID :many
SELECT * FROM wallets
WHERE user_id = $1;

-- name: GetWalletByID :one
SELECT * FROM wallets
WHERE id = $1 LIMIT 1;

-- name: GetWalletByAddress :one
SELECT * FROM wallets
WHERE address = $1 LIMIT 1;

-- name: UpdateWallet :one
UPDATE wallets SET
    address = $2,
    encrypted_private_key = $3,
    name = $4,
    status = $5,
    updated_at = $6
WHERE id = $1 RETURNING *;
