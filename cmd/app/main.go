package main

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"user-service/internal/config"
	"user-service/internal/db"
	"user-service/internal/kafka"
	"user-service/internal/metrics"
	"user-service/internal/repo/postgres"
	"user-service/internal/service"
	transport "user-service/internal/transport/http"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := db.NewPostgresPool(ctx, cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	producer := kafka.NewProducer(cfg.Kafka.BrokerAddr, cfg.Kafka.Topic)
	defer producer.Close()

	consumer := kafka.NewConsumer(cfg.Kafka.BrokerAddr, cfg.Kafka.Topic, cfg.Kafka.ConsumerGroup)
	go consumer.Start(ctx)
	defer consumer.Close()

	userRepo := postgres.NewUserRepo(pool)
	userService := service.NewUserService(userRepo, producer)
	handler := transport.NewHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/user", metrics.Middleware(handler.GetUser))
	mux.HandleFunc("/users", metrics.Middleware(handler.GetUsers))
	mux.HandleFunc("/user/create", metrics.Middleware(handler.CreateUser))
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
