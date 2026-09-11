package providers

import (
	"context"
	"log/slog"
	"time"

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
	slog.Info("[-]:Sending sms.....")
	time.Sleep(5 * time.Second)
	slog.Info("[OK]:Sending sms success.")
	return nil
}