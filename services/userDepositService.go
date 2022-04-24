package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"vendingmachine/config"
	"vendingmachine/domains"
	"vendingmachine/repositories"
)

type UserDepositService struct {
}

var (
	MakeUserDeposit = func(userDeposit *domains.UserDeposit) (*domains.UserDeposit, *config.HttpError) {
		return repositories.MakeUserDeposit(userDeposit)
	}
	ResetUserDeposit = func(id string) *config.HttpError {
		return repositories.DeleteUserDeposit(id)
	}
)

func NewUserDepositService() *UserDepositService {
	return &UserDepositService{}
}

func MakeDeposit(w http.ResponseWriter, r *http.Request) {

	var newDeposit domains.UserDeposit

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

	err = json.Unmarshal(reqBody, &newDeposit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadRequestError("Invalid Body \n " + err.Error()))

		return
	}

	if newDeposit.Deposit%5 != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(config.BadDepositError(string(rune(newDeposit.Deposit))))

		return
	}

	newDeposit.UserId = r.Header.Values("user_id")[0]

	ud, httpErr := MakeUserDeposit(&newDeposit)

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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			config.BadRequestError("user not found on the request"))
		return
	}
	userId := r.Header.Values("user_id")[0]

	httpErr := ResetUserDeposit(userId)

	if httpErr != nil {
		w.WriteHeader(httpErr.Code)
		json.NewEncoder(w).Encode(config.DataAccessLayerError(httpErr.Error))

		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w)
}
