-- name: ListPasswords :many
SELECT 
    id,
    user_id,
    service,
    url,
    notes,
    icon,
    created_at,
    updated_at
FROM passwords
WHERE user_id = $1;

-- name: GetPasswordByID :one
SELECT * FROM passwords
WHERE id = $1
    AND user_id = $2
LIMIT 1;

-- name: CreatePassword :one
INSERT INTO passwords (
    user_id,
    service,
    username,
    password,
    url,
    notes,
    icon
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, service, url, notes, icon, created_at, updated_at;

-- name: UpdatePassword :one
UPDATE passwords
SET 
    service = $3,
    username = $4,
    password = $5,
    url = $6,
    notes = $7,
    icon = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
RETURNING id, user_id, service, url, notes, icon, created_at, updated_at;

-- name: DeletePassword :one
DELETE FROM passwords
WHERE id = $1
    AND user_id = $2
RETURNING id, user_id, service, url, notes, icon, created_at, updated_at;
