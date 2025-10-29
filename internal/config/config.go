package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Kafka struct {
    BootstrapServersValue    string `validate:"required" env:"BOOTSTRAP_SERVERS_VALUE"`
    VkMessageTopicName       string `validate:"required" env:"TOPIC_NAME_VK_REQUEST_MESSAGE"`
    ResponseMessageTopicName string `validate:"required" env:"TOPIC_NAME_VK_RESPONSE_PREPARER"`
    ResponseMessageGroupID   string `validate:"required" env:"GROUP_ID_VK_RESPONSE_PREPARER"`
    ClientID                 string `env:"CLIENT_ID_VK_RESPONSE_PREPARER"`
    SaslUsername             string `env:"SASL_USERNAME"`
    SaslPassword             string `env:"SASL_PASSWORD"`
}

type Config struct {
    Kafka  Kafka  `envPrefix:"KAFKA_"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}
	v := validator.New()
	if err := v.Struct(c); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}
	return &c, nil
}
