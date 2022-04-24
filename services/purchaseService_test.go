package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"vendingmachine/config"
	"vendingmachine/domains"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name                   string
	userID                 string
	testID                 string
	mockGetProductById     func()
	mockGetUserDepositById func()
	mockRunFinishPurchase  func()
}

func TestPurchaseProductInsufficiency(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/buy/1/14", nil)
	r = mux.SetURLVars(r, map[string]string{"productId": "1", "quantity": "14"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")

	tempGetProductByID := GetProductByID

	testPurchaseProductInsufficiency := &test{

		name:   "Product Insufficiency Error Expected",
		userID: "seller",
		testID: "ProductInsufficiency",
		mockGetProductById: func() {
			GetProductByID = func(id int) *domains.Product {
				return &domains.Product{
					ID:              1,
					AmountAvailable: 10,
					Cost:            200,
					ProductName:     "Test Coffee",
					SellerId:        "seller-user-test",
				}
			}
		},
	}

	testPurchaseProductInsufficiency.mockGetProductById()
	Purchase(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	GetProductByID = tempGetProductByID

}

func TestPurchaseProductNotFound(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/buy/1/14", nil)
	r = mux.SetURLVars(r, map[string]string{"productId": "1", "quantity": "14"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")

	tempGetProductByID := GetProductByID

	testPurchaseProductNotFound := &test{

		name:   "Product Insufficiency Error Expected",
		userID: "seller",
		testID: "ProductInsufficiency",
		mockGetProductById: func() {
			GetProductByID = func(id int) *domains.Product {
				return nil
			}
		},
	}

	testPurchaseProductNotFound.mockGetProductById()
	Purchase(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	GetProductByID = tempGetProductByID

}

func TestUserDepositNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/buy/1/14", nil)
	r = mux.SetURLVars(r, map[string]string{"productId": "1", "quantity": "14"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")

	tempGetProductByID := GetProductByID

	testPurchaseUserDepositNotFound := &test{

		name:   "Product Insufficiency Error Expected",
		userID: "seller",
		testID: "ProductInsufficiency",
		mockGetProductById: func() {
			GetProductByID = func(id int) *domains.Product {
				return &domains.Product{
					ID:              1,
					AmountAvailable: 18,
					Cost:            200,
					ProductName:     "Test Coffee",
					SellerId:        "seller-user-test",
				}
			}
		},
		mockGetUserDepositById: func() {
			GetUserDepositByID = func(userId string) (*domains.UserDeposit, *config.HttpError) {
				return nil, config.DataAccessLayerError("User Deposit Not Found")
			}
		},
	}

	testPurchaseUserDepositNotFound.mockGetProductById()
	testPurchaseUserDepositNotFound.mockGetUserDepositById()
	Purchase(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	GetProductByID = tempGetProductByID

}

func TestUserDepositInsertMoreCoins(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/buy/1/14", nil)
	r = mux.SetURLVars(r, map[string]string{"productId": "1", "quantity": "14"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer")

	tempGetProductByID := GetProductByID

	testPurchaseUserDepositNotFound := &test{

		name:   "Insert more Coins Error Expected",
		userID: "buyer",
		testID: "ProductInsufficiency",
		mockGetProductById: func() {
			GetProductByID = func(id int) *domains.Product {
				return &domains.Product{
					ID:              1,
					AmountAvailable: 18,
					Cost:            200,
					ProductName:     "Test Coffee",
					SellerId:        "seller-user-test",
				}
			}
		},
		mockGetUserDepositById: func() {
			GetUserDepositByID = func(userId string) (*domains.UserDeposit, *config.HttpError) {
				return &domains.UserDeposit{
					Balance: 150,
					UserId:  "buyer",
				}, nil
			}
		},
	}

	testPurchaseUserDepositNotFound.mockGetProductById()
	testPurchaseUserDepositNotFound.mockGetUserDepositById()
	Purchase(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	GetProductByID = tempGetProductByID

}

func TestNoErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/buy/1/14", nil)
	r = mux.SetURLVars(r, map[string]string{"productId": "1", "quantity": "1"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "buyer")

	tempGetProductByID := GetProductByID

	testNoErrors := &test{

		name:   "No Errors",
		userID: "buyer",
		testID: "NoErrors",
		mockGetProductById: func() {
			GetProductByID = func(id int) *domains.Product {
				return &domains.Product{
					ID:              1,
					AmountAvailable: 18,
					Cost:            200,
					ProductName:     "Test Coffee",
					SellerId:        "seller-user-test",
				}
			}
		},
		mockGetUserDepositById: func() {
			GetUserDepositByID = func(userId string) (*domains.UserDeposit, *config.HttpError) {
				return &domains.UserDeposit{
					Balance: 250,
					UserId:  "buyer",
				}, nil
			}
		},
		mockRunFinishPurchase: func() {
			RunFinishPurchase = func(product domains.Product, userId string) error {
				return nil
			}
		},
	}

	testNoErrors.mockGetProductById()
	testNoErrors.mockGetUserDepositById()
	testNoErrors.mockRunFinishPurchase()
	Purchase(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	GetProductByID = tempGetProductByID

}
