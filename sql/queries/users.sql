-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(), -- generate a new UUID for id
    NOW(),             -- current timestamp for created_at
    NOW(),             -- current timestamp for updated_at
    $1                 -- first parameter passed in for email
)
RETURNING *;
