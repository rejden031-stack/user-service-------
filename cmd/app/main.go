package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"user-service/internal/controller"
	"user-service/internal/db"
	"user-service/internal/kafka"
	"user-service/internal/metrics"
	"user-service/internal/repo/postgres"
	"user-service/internal/service"
)

func main() {
	ctx := context.Background()

	brokerAddr := os.Getenv("KAFKA_BROKER_ADDR")
	if brokerAddr == ""{
		brokerAddr = "localhost:9092"
	}

	pool, err := db.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	producer := kafka.NewProducer(brokerAddr, "user-events")
	defer producer.Close()

	consumer := kafka.NewConsumer(brokerAddr, "user-events", "user-service-group")
	go consumer.Start(ctx)
	defer consumer.Close()

	
	userRepo := postgres.NewUserRepo(pool)
	userService := service.NewUserService(userRepo, producer)
	handler := controller.NewHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/user", metrics.Middleware(handler.GetUser))
	mux.HandleFunc("/users", metrics.Middleware(handler.GetUsers))
	mux.HandleFunc("/user/create", metrics.Middleware(handler.CreateUser))
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
