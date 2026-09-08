package providers

import (
	"context"
	
	"github.com/MyintMyatt/notification-service/internal/models"
)

type EmailProvider struct {}

func NewEmailProvider() *EmailProvider {
	return &EmailProvider{}
}

func (e *EmailProvider) Channel() models.Channel {
	return models.ChannelEmail
}

func (e *EmailProvider) SendNotification(ctx context.Context, request *models.NotificationRequest) error {
	// Implement the logic to send email notification here
	return nil
}