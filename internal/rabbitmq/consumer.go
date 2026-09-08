package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/MyintMyatt/notification-service/internal/models"
	"github.com/MyintMyatt/notification-service/internal/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RetryHeader = "x-retry-count"
	OriginalQueue = "x-original-queue"
)

type ConsumerConfig struct {
	Queue string
	ChannelType models.Channel
	Workers int
	Prefetch int
}

type ConsumerManager struct {
	client *RabbitMQClient
	providers ports.ProviderRegistry
	maxRetries int
	log *slog.Logger

	mu sync.Mutex
	cancel context.CancelFunc
	running bool
	wg sync.WaitGroup
}

func NewConsumerManager(client *RabbitMQClient, providers ports.ProviderRegistry, maxRetries int, log *slog.Logger) *ConsumerManager {
	if log == nil {
		log = slog.Default()
	}

	return &ConsumerManager{
		client: client,
		providers: providers,
		maxRetries: maxRetries,
		log: log,
	}
}

func (cm *ConsumerManager) Start(ctx context.Context, configs []ConsumerConfig) error {
	cm.mu.Lock()
	if cm.running {
		cm.mu.Unlock()
		return fmt.Errorf("consumer manager already running")
	}

	consumerCtx, cancel := context.WithCancel(ctx)
	cm.cancel = cancel
	cm.running = true
	cm.mu.Unlock()

	for _, config := range configs {
		cfg := config
		cm.wg.Add(1)
		go func() {
			defer cm.wg.Done()
			if err := cm.startWorkerPool(consumerCtx, cfg); err != nil {
				cm.log.Error("failed to start worker pool", "queue", cfg.Queue, "error", err)
			}
		} ()
	}

	cm.log.Info("ConsumerManager started", "queues", len(configs))
	return nil
}

func (cm *ConsumerManager) startWorkerPool(ctx context.Context, cfg ConsumerConfig) error {
	ch, err := cm.client.connection.Channel() // Open a new channel for this worker pool
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	defer ch.Close()

	prefetch := cfg.Prefetch
	if prefetch <= 0 {
		prefetch = cfg.Workers * 2
	}

	if err := ch.Qos(prefetch, 0, false); err != nil {
		return err
	}

	deliveries, err := ch.Consume(cfg.Queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	jobs := make(chan amqp.Delivery, prefetch)
	var wg sync.WaitGroup
	
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			cm.worker(ctx, workerID, jobs, cfg.ChannelType, cfg.Queue)
		}(i)
	}

	cm.log.Info("worker pool started",
		"queue", cfg.Queue,
		"workers", cfg.Workers,
		"channel", cfg.ChannelType,
	)

	// dispatcher
	go func() {
		defer close(jobs)
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					return
				}
				select {
				case jobs <- d:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	wg.Wait()
	return nil

}

func (cm *ConsumerManager) worker(ctx context.Context, workerID int, jobs <-chan amqp.Delivery, channelType models.Channel, queueName string) {
	for {
		select {
		case <-ctx.Done():
			cm.log.Info("worker shutting down", "workerID", workerID)
			return
		case delivery, ok := <-jobs:
			if !ok {
				cm.log.Info("jobs channel closed, worker exiting", "workerID", workerID)
				return
			}

			cm.processMessage(ctx, workerID, delivery, channelType, queueName)
		}
	}
}	

func (m *ConsumerManager) processMessage(
	ctx context.Context,
	workerID int,
	d amqp.Delivery,
	channelType models.Channel,
	queueName string,
) {
	var notif models.NotificationRequest // or models.NotificationRequest
	if err := json.Unmarshal(d.Body, &notif); err != nil {
		m.log.Error("malformed message", "worker", workerID, "error", err)
		_ = d.Nack(false, false) // → DLQ
		return
	}

	if err := notif.Validate(); err != nil {
		m.log.Error("validation failed", "id", notif.ID, "error", err)
		_ = d.Nack(false, false)
		return
	}

	provider, exists := m.providers[channelType]
	if !exists {
		m.log.Error("no provider", "channel", channelType)
		_ = d.Nack(false, false)
		return
	}

	// --- Idempotency (strongly recommended) ---
	// if alreadyProcessed := redis.SetNX(ctx, "notif:"+notif.ID, "1", 24*time.Hour); !alreadyProcessed {
	//     _ = d.Ack(false)
	//     return
	// }

	err := provider.SendNotification(ctx, &notif) // adapt to your interface
	if err == nil {
		m.log.Info("sent successfully",
			"id", notif.ID,
			"channel", channelType,
			"queue", queueName,
			"worker", workerID,
		)
		_ = d.Ack(false)
		return
	}

	retryCount := getRetryCount(d.Headers)
	if retryCount >= m.maxRetries {
		m.log.Error("max retries reached → DLQ",
			"id", notif.ID,
			"retries", retryCount,
			"error", err,
		)
		_ = d.Nack(false, false) // goes to your DLX
		return
	}

	retryCount++
	m.log.Warn("will retry",
		"id", notif.ID,
		"attempt", retryCount,
		"max", m.maxRetries,
		"error", err,
	)

	if err := m.publishToRetryQueue(d, channelType, retryCount, queueName); err != nil {
		m.log.Error("failed to publish to retry queue", "id", notif.ID, "error", err)
		_ = d.Nack(false, true) // last resort requeue
		return
	}

	_ = d.Ack(false) // remove from main queue
}

func getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	v, ok := headers[RetryHeader]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	default:
		return 0
	}
}

func (m *ConsumerManager) publishToRetryQueue(
	d amqp.Delivery,
	channelType models.Channel,
	retryCount int,
	originalQueue string,
) error {
	ch := m.client.channel // or open a short-lived channel
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}

	retryQueue := retryQueueName(channelType)
	if retryQueue == "" {
		return fmt.Errorf("no retry queue for %s", channelType)
	}

	headers := amqp.Table{}
	for k, v := range d.Headers {
		headers[k] = v
	}
	headers[RetryHeader] = retryCount
	headers[OriginalQueue] = originalQueue

	return ch.PublishWithContext(
		context.Background(),
		"", // default exchange
		retryQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  d.ContentType,
			Body:         d.Body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

func retryQueueName(ch models.Channel) string {
	switch ch {
	case models.ChannelEmail:
		return "q.retry.email"
	case models.ChannelSMS:
		return "q.retry.sms"
	case models.ChannelPush:
		return "q.retry.fcm"
	default:
		return ""
	}
}