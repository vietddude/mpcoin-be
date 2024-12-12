-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET
    email = $2,
    password_hash = $3,
    status = $4,
    updated_at = $5
WHERE id = $1 RETURNING *;
