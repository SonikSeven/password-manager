-- name: ListPasswords :many
SELECT * FROM passwords
WHERE user_id = $1;

-- name: GetPasswordByID :one
SELECT * FROM passwords
WHERE id = $1
    AND user_id = $2
LIMIT 1;

-- name: CreatePassword :one
INSERT INTO passwords (
    user_id,
    name,
    password
) VALUES (
    $1, $2, $3
)
RETURNING id, user_id, created_at, updated_at;

-- name: UpdatePassword :one
UPDATE passwords
SET 
    name = $3,
    password = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
RETURNING id, user_id, name, created_at, updated_at;

-- name: DeletePassword :one
DELETE FROM passwords
WHERE id = $1
    AND user_id = $2
RETURNING id, user_id, created_at, updated_at;
