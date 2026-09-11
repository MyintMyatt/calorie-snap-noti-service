package providers

import (
	"context"
	"log/slog"
	"time"

	"github.com/MyintMyatt/notification-service/internal/models"
)

type FcmProvider struct {}

func NewFcmProvider() *FcmProvider {
	return &FcmProvider{}
}

func (f *FcmProvider) Channel() models.Channel {
	return models.ChannelPush
}

func (f *FcmProvider) SendNotification(ctx context.Context, request *models.NotificationRequest) error {
	// Implement the logic to send FCM notification here
	slog.Info("[-]:Sending push noti.....")
	time.Sleep(5 * time.Second)
	slog.Info("[OK]:Sending push noti success.")
	return nil
}