package rabbitmq

import (
	"github.com/MyintMyatt/notification-service/internal/models"
)

func Config() []ConsumerConfig {
	return []ConsumerConfig{
		// Security Queues
		{Queue: SecurityEmailQueue, ChannelType: models.ChannelEmail, Workers: 8, Prefetch: 8},
		{Queue: SecuritySmsQueue, ChannelType: models.ChannelSMS, Workers: 6, Prefetch: 6},
		{Queue: SecurityFCMQueue, ChannelType: models.ChannelPush, Workers: 6, Prefetch: 6},

		// Retry Queues
		{Queue: RetryEmailQueue, ChannelType: models.ChannelEmail, Workers: 4, Prefetch: 4},
		{Queue: RetrySmsQueue, ChannelType: models.ChannelSMS, Workers: 4, Prefetch: 4},
		{Queue: RetryFCMQueue, ChannelType: models.ChannelPush, Workers: 4, Prefetch: 4},

		// User Queues
		{Queue: UserEmailQueue, ChannelType: models.ChannelEmail, Workers: 5, Prefetch: 8},
		{Queue: UserSmsQueue, ChannelType: models.ChannelSMS, Workers: 5, Prefetch: 6},
		{Queue: UserFCMQueue, ChannelType: models.ChannelPush, Workers: 8, Prefetch: 6},

		// Marketing Queues
		{Queue: AdvertisingEmailQueue, ChannelType: models.ChannelEmail, Workers: 3, Prefetch: 3},
		{Queue: AdvertisingSmsQueue, ChannelType: models.ChannelSMS, Workers: 3, Prefetch: 3},
		{Queue: AdvertisingFCMQueue, ChannelType: models.ChannelPush, Workers: 3, Prefetch: 3},

	}
}