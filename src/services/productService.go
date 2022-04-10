package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"vendingmachine/src/config"
	"vendingmachine/src/domains"
	"vendingmachine/src/repositories"

	"github.com/gorilla/mux"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	var newProduct domains.Product

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Problems reading your Request Body")
	}

	if json.Unmarshal(reqBody, &newProduct); int32(newProduct.Cost)%5 != 0 {
		//check if cost is multiple of 5
		fmt.Println(newProduct)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadProductCostError(string(rune(newProduct.Cost))))

		return

	}

	pd, httpErr := repositories.SaveProduct(&newProduct)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.UnauthorizedError())

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&pd)

}

func GetOneProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			config.BadRequestError("Id must be an integer"))
		return
	}

	product := repositories.FindProductById(id)

	if product == nil {
		w.WriteHeader(404)

		json.NewEncoder(w).Encode(config.NotFoundError())

		return
	}

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

	products := repositories.FindAllProducts()

	json.NewEncoder(w).Encode(&products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var updProduct domains.Product

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Problems reading your Request Body")
	}

	json.Unmarshal(reqBody, &updProduct)

	updProduct.SellerId = r.Header.Values("user_id")[0]

	pd, httpErr := repositories.UpdateProduct(&updProduct)

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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			config.BadRequestError("Id must be an integer"))
		return
	}

	httpErr := repositories.DeleteProduct(id)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w)
}
