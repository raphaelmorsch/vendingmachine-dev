package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"vendingmachine/src/domains"

	"github.com/gorilla/mux"
)

func TestPurchase(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/buy/1/14", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	m := mux.NewRouter()
	m.HandleFunc("/buy/{productId}/{quantity}", Purchase)
	m.ServeHTTP(w, r)

	tests := []struct {
		name               string
		userID             string
		mockGetProductById func()
	}{
		{
			name:   "All success",
			userID: "seller",
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tc.mockGetProductById()
			Purchase(w, r)
		})
	}

}
