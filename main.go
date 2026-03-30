package main

import (
	"log"
	"net/http"
	"os"

	"payment-gateway-test-manjo/database"
	"payment-gateway-test-manjo/handler"
	"payment-gateway-test-manjo/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	database.InitDB()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
    	AllowedOrigins:   []string{"*"},
    	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-SIGNATURE"},
    	ExposedHeaders:   []string{"Link"},
    	AllowCredentials: true,
	}))

	r.With(
		middleware.SignatureMiddleware(middleware.GenerateQRPayload),
	).Post("/api/v1/qr/generate", handler.GenerateQR)

	r.With(
		middleware.SignatureMiddleware(middleware.PaymentPayload),
	).Post("/api/v1/qr/payment", handler.PaymentCallback)

	r.Post("/api/v1/signature", handler.GenerateSignatureHandler)

	r.Get("/api/v1/transactions", handler.GetTransactions(database.DB))

	r.Get("/api/v1/tracker/{referenceNo}", handler.GetTransactionByReference(database.DB))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	http.ListenAndServe(":"+port, r)
}