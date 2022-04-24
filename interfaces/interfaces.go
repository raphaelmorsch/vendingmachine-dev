package interfaces

import (
	"net/http"
	"vendingmachine/config"
	"vendingmachine/domains"
)

type UserDepositRepositoryInterface interface {
	FindUserDeposit(userId string) (*domains.UserDeposit, *config.HttpError)

	MakeUserDeposit(userDeposit *domains.UserDeposit) (*domains.UserDeposit, *config.HttpError)

	DeleteUserDeposit(id string) *config.HttpError
}

type UserDepositServiceInterface interface {
	MakeDeposit(w http.ResponseWriter, r *http.Request)

	DeleteUserDeposit(w http.ResponseWriter, r *http.Request)
}
