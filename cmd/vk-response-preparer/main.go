package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"vk-response-preparer/internal/config"
	"vk-response-preparer/internal/logger"
	"vk-response-preparer/internal/messaging"
	"vk-response-preparer/internal/processor"
)

func main() {
	lg, cleanup := logger.NewZapLogger()
	defer cleanup()

	lg.Info("🚀 Starting vk-response-preparer…")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		lg.Error("❌ Failed to load config: %v", err)
		os.Exit(1)
	}

	producer, err := messaging.NewKafkaProducer(messaging.Option{
		Logger:       lg,
		Broker:       cfg.Kafka.BootstrapServersValue,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
		ClientID:     cfg.Kafka.ClientID,
	})
	if err != nil {
		lg.Error("❌ Failed to create producer: %v", err)
		os.Exit(1)
	}
	defer producer.Close()

	vkMessagePreparer := processor.NewVkMessagePreparer(cfg.Kafka.VkMessageTopicName, producer)

	consumer, err := messaging.NewKafkaConsumer(messaging.ConsumerOption{
		Logger:       lg,
		Broker:       cfg.Kafka.BootstrapServersValue,
		GroupID:      cfg.Kafka.ResponseMessageGroupID,
		Topics:       []string{cfg.Kafka.ResponseMessageTopicName},
		Handler:      vkMessagePreparer,
		SaslUsername: cfg.Kafka.SaslUsername,
		SaslPassword: cfg.Kafka.SaslPassword,
		ClientID:     cfg.Kafka.ClientID,
	})
	if err != nil {
		lg.Error("❌ Failed to create consumer: %v", err)
		os.Exit(1)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	if err := consumer.Start(ctx); err != nil {
		lg.Error("❌ Consumer error: %v", err)
		os.Exit(1)
	}
}
