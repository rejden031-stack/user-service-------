package main

import (
	"context"
	"log"
	"net/http"
	"user-service/internal/db"
	myhttp "user-service/internal/http"
	"user-service/internal/kafka"
	"user-service/internal/service"
)

func main() {
	ctx := context.Background()

	pool, err := db.NewPostgresPool(ctx)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	producer := kafka.NewProducer("localhost:9092", "user-events")
	defer producer.Close()

	// Запускаем consumer в отдельной горутине
	consumer := kafka.NewConsumer("localhost:9092", "user-events", "user-service-group")
	go consumer.Start(ctx)
	defer consumer.Close()

	userService := service.NewUserService(pool)
	handler := myhttp.NewHandler(userService, producer)

	http.HandleFunc("/user", handler.GetUser)
	http.HandleFunc("/users", handler.GetUsers)
	http.HandleFunc("/user/create", handler.CreateUser)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
