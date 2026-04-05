package valueobjects

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSuccess NotificationStatus = "success"
	NotificationStatusFailed  NotificationStatus = "failed"
)
