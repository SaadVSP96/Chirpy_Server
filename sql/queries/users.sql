-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    FALSE
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red;


-- name: DeletAllUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    hashed_password = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET 
    is_chirpy_red = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
