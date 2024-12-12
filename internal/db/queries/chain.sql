-- name: GetChainByID :one
SELECT * FROM chains WHERE id = $1;

-- name: GetChainByChainID :one
SELECT * FROM chains WHERE chain_id = $1;

-- name: GetChains :many
SELECT * FROM chains;

