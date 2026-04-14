package main

import (
	"log"
	"net/http"
	myhttp "user-service/internal/http"
)

func main() {
	handler := myhttp.NewHandler()
	http.HandleFunc("/user", handler.GetUser)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
