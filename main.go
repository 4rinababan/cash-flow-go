package main

import (
	"log"
	"net/http"

	db "cash-flow-go/database"
	"cash-flow-go/handlers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "cash-flow-go/docs" // ini penting buat Swagger
	// "github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
)

func main() {
	db.Init() // connect DB + migrate

	r := mux.NewRouter()

	r.HandleFunc("/api/transactions", handlers.CreateTransaction).Methods("POST")
	r.HandleFunc("/api/transactions", handlers.GetTransactions).Methods("GET")
	r.HandleFunc("/api/transactions/{id}", handlers.DeleteTransaction).Methods("DELETE")
	r.HandleFunc("/api/transactions/top5", handlers.GetTop5Transactions).Methods("GET")

	r.HandleFunc("/api/dashboard", handlers.GetDashboard).Methods("GET")
	r.HandleFunc("/api/dashboard/bar", handlers.GetBarChart).Methods("GET")
	r.HandleFunc("/api/dashboard/donut", handlers.GetDonutChart).Methods("GET")
	r.HandleFunc("/api/dashboard/monthly-bar", handlers.GetMonthlyBarChart).Methods("GET")

	// Swagger endpoint
	// Swagger endpoint (pastikan pakai handler, bukan "value")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running at :8889")
	log.Fatal(http.ListenAndServe(":8889", r))
}
