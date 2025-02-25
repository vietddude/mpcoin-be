-- name: CreateTransaction :one
INSERT INTO transactions (chain_id, from_address, to_address, tx_hash, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: GetTransactionsByWalletAddress :many
SELECT * FROM transactions 
WHERE (from_address = $1 OR to_address = $1) 
AND ($2::int IS NULL OR chain_id = $2)
ORDER BY created_at DESC
LIMIT $3
OFFSET $4;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1;

-- name: GetTransactionCount :one
SELECT COUNT(*) 
FROM transactions 
WHERE (from_address = $1 OR to_address = $1)
  AND ($2::int IS NULL OR chain_id = $2);