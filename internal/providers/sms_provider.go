package providers

import (
	"context"
	
	"github.com/MyintMyatt/notification-service/internal/models"
)
type SmsProvider struct {}

func NewSmsProvider() *SmsProvider {
	return &SmsProvider{}
}

func (s *SmsProvider) Channel() models.Channel {
	return models.ChannelSMS
}

func (s *SmsProvider) SendNotification(ctx context.Context, request *models.NotificationRequest) error {
	// Implement the logic to send SMS notification here
	return nil
}