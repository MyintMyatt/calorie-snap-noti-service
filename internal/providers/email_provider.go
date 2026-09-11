package providers

import (
	"context"
	"log/slog"
	"time"

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
	slog.Info("[-]:Sending mail.....")
	time.Sleep(5 * time.Second)
	slog.Info("[OK]:Sending mail success.")
	return nil
}