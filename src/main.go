package main

import (
	"net/http"
	"vendingmachine/src/config"
	"vendingmachine/src/controllers"
	"vendingmachine/src/services"

	"github.com/gorilla/mux"

	"log"
)

func main() {

	run()
}

func run() {
	config.DBConnect()

	services.InitializeOauthServer()

	router := mux.NewRouter().StrictSlash(true)

	router.Use(commonMiddleware)

	registerRoutes(router)

	log.Fatal(http.ListenAndServe(":8083", router))

}

func registerRoutes(router *mux.Router) {

	registerControllerRoutes(controllers.ProductController{}, router)
	registerControllerRoutes(controllers.UserController{}, router)
	registerControllerRoutes(controllers.SwaggerController{}, router)
	registerControllerRoutes(controllers.UserDepositController{}, router)
	registerControllerRoutes(controllers.PurchaseController{}, router)

}

func registerControllerRoutes(controller controllers.Controller, router *mux.Router) {

	controller.RegisterRoutes(router)
}

func commonMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
