package controllers

import (
	"net/http"
	"vendingmachine/src/services"

	"github.com/gorilla/mux"
)

type UserController struct {
}

func (t UserController) RegisterRoutes(router *mux.Router) {

	router.Handle("/user", http.HandlerFunc(services.AddNewUser)).Methods("POST")

}
