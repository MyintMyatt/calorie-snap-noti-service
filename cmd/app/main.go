package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MyintMyatt/notification-service/internal/config"
	"github.com/MyintMyatt/notification-service/internal/ports"
	"github.com/MyintMyatt/notification-service/internal/providers"
	"github.com/MyintMyatt/notification-service/internal/rabbitmq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("[-]:Connecting Rabbitmq........")
	rabbitClient, err := rabbitmq.NewRabbitMQClient(cfg.RabbitMQUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitClient.Close()

	fmt.Println("[OK]:Connected Rabbitmq.")

	emailProvider := providers.NewEmailProvider()
	fcmProvider := providers.NewFcmProvider()
	smsProvder := providers.NewSmsProvider()


	providers := ports.ProviderRegistry{
		emailProvider.Channel(): emailProvider,
		fcmProvider.Channel() : fcmProvider,
		smsProvder.Channel() : smsProvder,
	}

	cm := rabbitmq.NewConsumerManager(rabbitClient, providers, 3, slog.Default());

	// Root context that can be cancelled on signal
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, // Ctrl + C
		syscall.SIGTERM, // docker stop / k8s terminate
	)

	defer stop()

	if err := cm.Start(ctx, rabbitmq.Config()); err != nil {
		log.Fatal("failed to start consumer manager", "error", err)
		os.Exit(1)
	}

	slog.Info("[OK]:Notification Service Started")
	<- ctx.Done()
	slog.Info("shutdown signal received, stopping consumers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 25 * time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		cm.Stop()
		close(done)
	}()


	select {
	case <-done:
		slog.Info("graceful shutdown completed")
	case <-shutdownCtx.Done():
		slog.Warn("shutdown timed out, forcing exit")
	}
}
