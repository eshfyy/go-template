-- name: CreateUser :exec
INSERT INTO users (id, name, surname, telegram_id, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateUser :exec
UPDATE users SET name = $1, surname = $2, telegram_id = $3 WHERE id = $4;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
