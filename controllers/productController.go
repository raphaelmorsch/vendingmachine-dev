package controllers

import (
	"net/http"
	"vendingmachine/services"

	"github.com/gorilla/mux"
)

type ProductController struct{}

func (t ProductController) RegisterRoutes(router *mux.Router) {

	router.Handle("/product", services.Protect(http.HandlerFunc(services.CreateProduct), []string{"seller"})).Methods("POST")

	router.Handle("/product/{id}", services.Protect(http.HandlerFunc(services.GetOneProduct), []string{})).Methods("GET")

	router.Handle("/products", services.Protect(http.HandlerFunc(services.AllProducts), []string{})).Methods("GET")

	router.Handle("/product", services.Protect(http.HandlerFunc(services.UpdateProduct), []string{"seller"})).Methods("PUT")

	router.Handle("/product/{id}", services.Protect(http.HandlerFunc(services.DeleteProduct), []string{})).Methods("DELETE")
}
