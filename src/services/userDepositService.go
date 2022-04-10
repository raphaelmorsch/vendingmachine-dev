package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"vendingmachine/src/config"
	"vendingmachine/src/domains"
	"vendingmachine/src/repositories"
)

type UserDepositService struct {
}

func NewUserDepositService() *UserDepositService {
	return &UserDepositService{}
}

func MakeDeposit(w http.ResponseWriter, r *http.Request) {

	var newDeposit domains.UserDeposit

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Problems reading your Request Body")
	}

	json.Unmarshal(reqBody, &newDeposit)

	newDeposit.UserId = r.Header.Values("user_id")[0]

	if newDeposit.Deposit%5 != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadDepositError(string(rune(newDeposit.Deposit))))

		return
	}

	ud, httpErr := repositories.MakeUserDeposit(&newDeposit)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&ud)

}

func DeleteUserDeposit(w http.ResponseWriter, r *http.Request) {

	if len(r.Header.Values("user_id")) == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(
			config.BadRequestError("user not found on the request"))
		return
	}
	userId := r.Header.Values("user_id")[0]

	httpErr := repositories.DeleteUserDeposit(userId)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w)
}
