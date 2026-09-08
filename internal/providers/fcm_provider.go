package providers

import (
	"context"
	
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
	return nil
}