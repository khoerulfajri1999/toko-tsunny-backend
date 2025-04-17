package routes

import (
	"go-jwt/controllers"
	"go-jwt/middlewere"

	"github.com/gorilla/mux"
)

func CategoryRoutes(r *mux.Router) {
	router := r.PathPrefix("/category").Subrouter()

	router.Use(middlewere.Auth)

	router.HandleFunc("", controllers.CreateCategory).Methods("POST")
	router.HandleFunc("", controllers.GetAllCategory).Methods("GET")
	router.HandleFunc("/{id}", controllers.GetCategoryById).Methods("GET")
	router.HandleFunc("/{id}", controllers.UpdateCategoryById).Methods("PUT")
	router.HandleFunc("/{id}", controllers.DeleteCategoryById).Methods("DELETE")
}
