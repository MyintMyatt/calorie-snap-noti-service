package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MyintMyatt/notification-service/internal/config"
	"github.com/MyintMyatt/notification-service/internal/event"
	"github.com/MyintMyatt/notification-service/internal/rabbitmq"
)

func main() {
	welcome := event.WelcomeUser{
		UserId: "testing",
		Email:  "hello@gmail.com",
	}

	e := event.Event{
		ID:         "event-123",
		Type:       "user.welcome",
		Version:    1,
		Source:     "user-service",
		ClientTime: time.Now(),
		Data:       welcome,
	}

	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting Rabbitmq........")
	rabbitConn, err := rabbitmq.Connect(cfg.RabbitMQUrl)
	if err != nil {
		panic(err)
	}
	defer rabbitConn.Close()

	fmt.Println("Connected Rabbitmq........")
}
