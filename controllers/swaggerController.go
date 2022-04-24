package controllers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

type SwaggerController struct {
}

func (t SwaggerController) RegisterRoutes(router *mux.Router) {

	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)

}
