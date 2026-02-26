package models

// EmailNotificationTarget represents an email address that should receive notifications.
type EmailNotificationTarget struct {
	ID    int
	Email string
	Name  string
}
