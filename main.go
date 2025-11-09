package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	connectDB()

	r := chi.NewRouter()

	r.Post("/items", createItem)
	r.Get("/items", getItems)
	r.Get("/items/{id}", getItem)
	r.Put("/items/{id}", updateItem)
	r.Delete("/items/{id}", deleteItem)

	port := os.Getenv("PORT")
	log.Printf("🚀 Server running on port %s\n", port)
	http.ListenAndServe(":"+port, r)
}
