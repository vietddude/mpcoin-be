-- name: GetTokensByChain :many
SELECT * FROM tokens WHERE chain_id = $1;

-- name: GetTokenBySymbol :one
SELECT * FROM tokens WHERE chain_id = $1 AND symbol = $2;

-- name: GetTokenByContractAddress :one
SELECT * FROM tokens WHERE contract_address = $1;

-- name: GetTokenByID :one
SELECT * FROM tokens WHERE id = $1;