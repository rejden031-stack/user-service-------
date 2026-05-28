package main

import (
	"log"
	"net/http"
	"time"
	myhttp "user-service/internal/http"
)

func main() {
	ch := make(chan string)

	go func() {
		for i := 0; i <= 5; i++ {
			time.Sleep(1 * time.Second)

			ch <- "Фоновое сообщение"
		}
	}()

	go func() {
		for msg := range ch {
			log.Println("Получено:", msg)
		}
	}()

	handler := myhttp.NewHandler()
	http.HandleFunc("/user", handler.GetUser)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) //прикол
	log.Println("test change")
}
