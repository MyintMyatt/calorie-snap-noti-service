package ports

import (
	"context"

	"github.com/MyintMyatt/notification-service/internal/models"
)

type NotificationProvider interface {
	SendNotification(ctx context.Context, request *models.NotificationRequest) error
	Channel() models.Channel
}

type ProviderRegistry map[models.Channel]NotificationProvider