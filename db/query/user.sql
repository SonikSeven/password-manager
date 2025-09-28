-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4
) RETURNING id, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;
