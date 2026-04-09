-- name: CreateNotification :exec
INSERT INTO notifications (id, user_id, title, text, channel, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetNotificationByID :one
SELECT * FROM notifications WHERE id = $1;

-- name: ListNotificationsByUserID :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateNotificationStatus :exec
UPDATE notifications SET status = $1, sent_at = $2 WHERE id = $3;

-- name: ListFailedNotifications :many
SELECT * FROM notifications
WHERE status = 'failed' AND created_at >= now() - @since::interval
ORDER BY created_at ASC
LIMIT $1;

-- name: DeleteNotification :exec
DELETE FROM notifications WHERE id = $1;
