-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),  -- id
    NOW(),              -- created_at
    NOW(),              -- updated_at
    $1,                 -- body
    $2                  -- user_id
)
RETURNING *;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: GetChirpByID :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: ListChirps :many
SELECT *
FROM chirps
ORDER BY created_at DESC;

-- name: ListChirpsByUser :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;
