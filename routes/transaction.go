package routes

import (
	"go-jwt/controllers"
	"go-jwt/middlewere"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	router := r.PathPrefix("/transaction").Subrouter()

	router.Use(middlewere.Auth)

	router.HandleFunc("", controllers.CreateTransaction).Methods("POST")
	router.HandleFunc("", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/{id}", controllers.GetTransactionById).Methods("GET")
}
