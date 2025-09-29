-- name: ListPasswords :many
SELECT id, user_id, url, notes, icon, created_at, updated_at
FROM passwords
WHERE user_id = $1
  AND (
    @search::text IS NULL OR username ILIKE '%' || @search || '%'
    OR url ILIKE '%' || @search || '%'
    OR notes ILIKE '%' || @search || '%'
  )
  AND (
    @filter::text IS NULL
    OR url ~* ('(https?://)?([^/]*\.)?' || @filter)
  )
ORDER BY created_at DESC;

-- name: GetPasswordByID :one
SELECT * FROM passwords
WHERE id = $1
    AND user_id = $2
LIMIT 1;

-- name: CreatePassword :one
INSERT INTO passwords (
    user_id,
    username,
    password,
    url,
    notes,
    icon
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, user_id, url, notes, icon, created_at, updated_at;

-- name: UpdatePassword :one
UPDATE passwords
SET 
    username = $3,
    password = $4,
    url = $5,
    notes = $6,
    icon = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND user_id = $2
RETURNING id, user_id, url, notes, icon, created_at, updated_at;

-- name: DeletePassword :one
DELETE FROM passwords
WHERE id = $1
    AND user_id = $2
RETURNING id, user_id, url, notes, icon, created_at, updated_at;
