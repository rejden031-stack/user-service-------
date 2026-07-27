package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Database
	HTTPServer
	Kafka
}

type Database struct {
	DSN string `env:"DATABASE_URL,required"`
}

type HTTPServer struct {
	ReadTimeout  int `env:"HTTP_READ_TIMEOUT" envDefault:"5"`
	WriteTimeout int `env:"HTTP_WRITE_TIMEOUT" envDefault:"10"`
}

type Kafka struct {
	BrokerAddr    string `env:"KAFKA_BROKER_ADDR" envDefault:"localhost:9092"`
	Topic         string `env:"KAFKA_TOPIC" envDefault:"user-events"`
	ConsumerGroup string `env:"KAFKA_CONSUMER_GROUP" envDefault:"user-service-group"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
