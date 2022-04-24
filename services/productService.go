package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"vendingmachine/config"
	"vendingmachine/domains"
	"vendingmachine/repositories"

	"github.com/gorilla/mux"
)

var (
	QueryProductById = func(id int) *domains.Product {
		return repositories.FindProductById(id)
	}
	SaveProduct = func(product *domains.Product) (*domains.Product, *config.HttpError) {
		return repositories.SaveProduct(product)
	}

	UpdateThisProduct = func(product *domains.Product) (*domains.Product, *config.HttpError) {
		return repositories.UpdateProduct(product)
	}

	DeleteThisProduct = func(id int) *config.HttpError {
		return repositories.DeleteProduct(id)
	}

	GetAllProducts = func() []domains.Product {
		return repositories.FindAllProducts()
	}
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	var newProduct domains.Product

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Invalid Body \n " + err.Error()))

		return
	}
	if len(reqBody) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Empty Body"))

		return
	}

	if json.Unmarshal(reqBody, &newProduct); int32(newProduct.Cost)%5 != 0 {
		//check if cost is multiple of 5
		log.Printf("Invalid product cost %v", newProduct.Cost)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadProductCostError(string(rune(newProduct.Cost))))

		return

	}

	pd, httpErr := SaveProduct(&newProduct)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&pd)

}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]
	if len(idStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			config.BadRequestError("Please inform Product ID"))
		return

	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			config.BadRequestError("Id must be an integer"))
		return
	}

	product := QueryProductById(id)

	if product == nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(config.NotFoundError())

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&product)
}

// swagger:route GET /products buyer listProduct
// Get products list
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetProduct
func AllProducts(w http.ResponseWriter, r *http.Request) {

	products := GetAllProducts()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var updProduct domains.Product

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Invalid Body structure \n " + err.Error()))

		return
	}
	if len(reqBody) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Empty Body"))

		return
	}

	if err = json.Unmarshal(reqBody, &updProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Invalid Body value \n" + err.Error()))

		return
	}

	if int32(updProduct.Cost)%5 != 0 {
		//check if cost is multiple of 5
		log.Printf("Invalid product cost %v", updProduct.Cost)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadProductCostError(string(rune(updProduct.Cost))))

		return

	}

	if len(r.Header.Values("user_id")) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Invalid Header missing user_id"))

		return
	}

	updProduct.SellerId = r.Header.Values("user_id")[0]

	pd, httpErr := UpdateThisProduct(&updProduct)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(&pd)

}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]
	if len(idStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			config.BadRequestError("Please inform Product ID"))

		return

	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Id must be an integer"))

		return
	}

	httpErr := DeleteThisProduct(id)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w)
}
