package models

import "time"

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
)

type NotificationPriority uint8

const (
	LowPriority    NotificationPriority = 1
	NormalPriority NotificationPriority = 2
	HighPriority   NotificationPriority = 4
	UrgentPriority NotificationPriority = 5
)

type Locale string

const (
	LocaleEN Locale = "en"
	LocaleMM Locale = "mm"
)

type NotificationRequest struct {
	ID         string                 `json:"id"`
	Channel    Channel                `json:"channel"` 
	Recipient  string                 `json:"recipient"` // user id or email
	Priority   NotificationPriority   `json:"priority"`
	ClientTime time.Time              `json:"client_time"`
	Template   string                 `json:"template"`
	Locale     Locale                 `json:"locale"`
	Data       map[string]interface{} `json:"data,omitempty"`
}
