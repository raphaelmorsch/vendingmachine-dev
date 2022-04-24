package main

import (
	"net/http"
	"vendingmachine/config"
	"vendingmachine/controllers"
	"vendingmachine/services"

	"github.com/gorilla/mux"

	"log"
)

func main() {

	run()
}

func run() {
	log.Println("Connecting to Database")
	config.DBConnect()
	log.Println("Database connected")

	log.Println("Initializing Authentication Issuer")
	services.InitializeOauthServer()
	log.Println("Authentication Issuer OK")

	log.Println("Setting Realm Admin User")
	services.AddRealmAdminUser()
	log.Println("Realm Admin Added Successfully")

	log.Println("Setting API Client Secret")
	services.UpdateAPIClientSercret()
	log.Println("Client Secret Refreshed Successfully")

	router := mux.NewRouter().StrictSlash(true)

	router.Use(commonMiddleware)

	log.Println("Registering routes")
	registerRoutes(router)
	log.Println("Registering routes OK")

	log.Println("API up and running")
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
