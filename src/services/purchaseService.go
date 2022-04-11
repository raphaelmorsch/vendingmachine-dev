package services

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vendingmachine/src/config"
	"vendingmachine/src/domains"
	"vendingmachine/src/repositories"

	"github.com/gorilla/mux"
)

var (
	GetProductByID = func(id int) *domains.Product {
		return repositories.FindProductById(int(id))
	}
	GetUserDepositByID = func(userId string) (*domains.UserDeposit, *config.HttpError) {
		return repositories.FindUserDeposit(userId)
	}
	RunFinishPurchase = func(product domains.Product, userId string) error {
		return repositories.FinishPurchase(product, userId)
	}
)

func Purchase(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	productId, errProductId := strconv.Atoi(params["productId"])
	if errProductId != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Check if Product ID is valid"))

		return
	}

	quantity, errQtd := strconv.Atoi(params["quantity"])
	if errQtd != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Check if quantity is valid"))

		return
	}

	userId := r.Header.Values("user_id")[0]

	//check if there are products availabe
	product := GetProductByID(productId)
	if product != nil {
		if product.AmountAvailable < quantity {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(config.BusinessLayerError("Product insufficiency"))

			return

		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Product not found"))
		json.NewEncoder(w).Encode(config.BadRequestError("Product not found"))

		return
	}

	//check if there is Deposit enough available
	deposit, errDeposit := GetUserDepositByID(userId)
	if errDeposit != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errDeposit.Error)

		return
	}
	if deposit.Balance < (int(product.Cost) * quantity) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(config.BusinessLayerError("Please, insert more coins"))

		return
	}

	var purchase domains.Purchase
	purchase.ProductId = product.ID
	purchase.TotalSpent = (int(product.Cost) * quantity)

	//subtract the amount available for the product
	product.AmountAvailable = product.AmountAvailable - quantity

	//subtract the cost from user's deposit
	deposit.Balance = deposit.Balance - (int(product.Cost) * quantity)

	//check if there is any change and if so, split it into coins
	if deposit.Balance > 0 {
		purchase.Change = ReturnChange(deposit.Balance)
	}

	errFinish := RunFinishPurchase(*product, userId)
	if errFinish != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errFinish.Error())

		return
	}

	//return Purchase json
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&purchase)
}

func ReturnChange(changeValue int) map[int]int {

	coinMap := make(map[int]int)
	if changeValue%5 != 0 {
		panic("invalid change value")
	}
	availableCoins := []int{100, 50, 20, 10, 5}

	for _, coin := range availableCoins {
		if (changeValue / coin) > 0 {
			coinMap[coin] = (changeValue / coin)
		}
		changeValue = changeValue - ((changeValue / coin) * coin)
	}

	return coinMap

}
