package routes

import (
	"go-jwt/controllers"
	"go-jwt/middlewere"

	"github.com/gorilla/mux"
)

func ProductRoutes(r *mux.Router) {
	router := r.PathPrefix("/product").Subrouter()

	router.Use(middlewere.Auth)

	router.HandleFunc("", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("", controllers.GetAllProduct).Methods("GET")
	router.HandleFunc("/{id}", controllers.GetProductById).Methods("GET")
	router.HandleFunc("/{id}", controllers.UpdateProductById).Methods("PUT")
	router.HandleFunc("/{id}", controllers.DeleteProductById).Methods("DELETE")
}
