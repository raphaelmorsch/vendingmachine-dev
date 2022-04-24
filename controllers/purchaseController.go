package controllers

import (
	"net/http"
	"vendingmachine/services"

	"github.com/gorilla/mux"
)

type PurchaseController struct{}

func (t PurchaseController) RegisterRoutes(router *mux.Router) {

	router.Handle("/buy/{productId}/{quantity}", services.Protect(http.HandlerFunc(services.Purchase), []string{"buyer"})).Methods("POST")

}
