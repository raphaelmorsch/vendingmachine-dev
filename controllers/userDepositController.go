package controllers

import (
	"net/http"
	"vendingmachine/services"

	"github.com/gorilla/mux"
)

type UserDepositController struct{}

func (t UserDepositController) RegisterRoutes(router *mux.Router) {

	router.Handle("/deposit", services.Protect(http.HandlerFunc(services.MakeDeposit), []string{"buyer"})).Methods("POST")

	router.Handle("/reset", services.Protect(http.HandlerFunc(services.DeleteUserDeposit), []string{"buyer"})).Methods("DELETE")
}
