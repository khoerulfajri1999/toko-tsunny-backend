package routes

import (
	"go-jwt/controllers"
	"go-jwt/middlewere"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	router := r.PathPrefix("/user").Subrouter()

	router.Use(middlewere.Auth)

	router.HandleFunc("/me", controllers.Me).Methods("GET")
	router.HandleFunc("/{id}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("", controllers.GetAllUser).Methods("GET")
	router.HandleFunc("/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/{id}", controllers.DeleteUser).Methods("DELETE")
}