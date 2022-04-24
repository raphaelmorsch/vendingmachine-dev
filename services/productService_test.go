package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vendingmachine/config"
	"vendingmachine/domains"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type testProduct struct {
	name                  string
	userID                string
	testID                string
	mockGetProductById    func()
	mockSaveProduct       func()
	mockGetAllProducts    func()
	mockUpdateThisProduct func()
	mockDeleteThisProduct func()
}

//TestCreateProductEmptyBody
func TestCreateProductEmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/product", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	CreateProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestCreatProductBadProductCost
func TestCreateProductBadProductCost(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(`{
		"amountAvailable": 15,				
		"cost": 153,			
		"productName": "Double Expresso",
		"sellerId": 1
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	CreateProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestCreateProductDataAccessLayerError
func TestCreateProductDataAccessLayerError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(`{
		"amountAvailable": 15,				
		"cost": 100,			
		"productName": "Double Expresso",
		"sellerId": 1
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	testCreateProductDALE := testProduct{
		name:   "Test Create Product Data Access Error",
		userID: "seller-user-test",
		testID: "TestCreateProductDataAccessError",
		mockSaveProduct: func() {
			SaveProduct = func(product *domains.Product) (*domains.Product, *config.HttpError) {
				return nil, config.DataAccessLayerError("")
			}
		},
	}
	testCreateProductDALE.mockSaveProduct()
	CreateProduct(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

}

//TestCreateProductValid
func TestCreateProductValid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(`{
		"amountAvailable": 15,				
		"cost": 100,			
		"productName": "Double Expresso",
		"sellerId": 1
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	testCreateProductDALE := testProduct{
		name:   "Test Create Product Valid",
		userID: "seller-user-test",
		testID: "TestCreateProductValid",
		mockSaveProduct: func() {
			SaveProduct = func(product *domains.Product) (*domains.Product, *config.HttpError) {
				return &domains.Product{
					ID:              100,
					AmountAvailable: 15,
					Cost:            100,
					ProductName:     "Double Expresso",
					SellerId:        "1",
				}, nil
			}
		},
	}
	testCreateProductDALE.mockSaveProduct()
	CreateProduct(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

//TestGetOneProductNoProductID
func TestGetOneProductNoProductID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/product", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	GetOneProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestGetOneProductInvalidProductID
func TestGetOneProductInvalidProductID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/product/s", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "s"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	GetOneProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestGetOneProductNotFound
func TestGetOneProductNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/product/21", nil)
	r.Header.Set("user_id", "buyer-user-test")
	r = mux.SetURLVars(r, map[string]string{"id": "21"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testGetOneProductNotFound := testProduct{
		name:   "Test Get One Product Not Found",
		userID: "buyer-user-test",
		testID: "TestGetOneProductNotFound",
		mockGetProductById: func() {
			QueryProductById = func(id int) *domains.Product {
				return nil
			}
		},
	}
	testGetOneProductNotFound.mockGetProductById()
	GetOneProduct(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)

}

//TestGetOneProductValid
func TestGetOneProductValid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/product/21", nil)
	r.Header.Set("user_id", "buyer-user-test")
	r = mux.SetURLVars(r, map[string]string{"id": "21"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testGetOneProductNotFound := testProduct{
		name:   "Test Get One Product Not Found",
		userID: "buyer-user-test",
		testID: "TestGetOneProductNotFound",
		mockGetProductById: func() {
			QueryProductById = func(id int) *domains.Product {
				return &domains.Product{
					ID:              21,
					AmountAvailable: 10,
					Cost:            150,
					ProductName:     "Mochiato",
					SellerId:        "1",
				}
			}
		},
	}
	testGetOneProductNotFound.mockGetProductById()
	GetOneProduct(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

}

//TestFindAllProducts
func TestFindAllProducts(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.Header.Set("user_id", "buyer-user-test")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testGetAllProducts := testProduct{
		name:   "Test Get All Products",
		userID: "buyer-user-test",
		testID: "TestGetAllProducts",
		mockGetAllProducts: func() {
			GetAllProducts = func() []domains.Product {
				return []domains.Product{
					{
						ID:              21,
						AmountAvailable: 10,
						Cost:            150,
						ProductName:     "Mochiato",
						SellerId:        "1",
					},
					{
						ID:              22,
						AmountAvailable: 10,
						Cost:            170,
						ProductName:     "Expresso",
						SellerId:        "1",
					},
					{
						ID:              24,
						AmountAvailable: 10,
						Cost:            190,
						ProductName:     "Double Expresso",
						SellerId:        "1",
					},
				}
			}
		},
	}
	testGetAllProducts.mockGetAllProducts()
	AllProducts(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

}

//TestFindNoProducts
func TestFindNoProducts(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.Header.Set("user_id", "buyer-user-test")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testGetAllProducts := testProduct{
		name:   "Test Get All Products",
		userID: "buyer-user-test",
		testID: "TestGetAllProducts",
		mockGetAllProducts: func() {
			GetAllProducts = func() []domains.Product {
				return nil
			}
		},
	}
	testGetAllProducts.mockGetAllProducts()
	AllProducts(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

}

//TestUpdateProductInvalidProductCost
func TestUpdateProductInvalidProductCost(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/product", strings.NewReader(`{
		"amountAvailable": 15,				
		"cost": 153,			
		"productName": "Double Expresso",
		"sellerId": "1"
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	UpdateProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestUpdateProductEmptyBody
func TestUpdateProductEmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/product", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	UpdateProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestUpdateProductInvalidBodyValue
func TestUpdateProductInvalidBodyValue(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/product", strings.NewReader(`{
		"amountAvailable": "invalidAmount",				
		"cost": 150,			
		"productName": "Double Expresso",
		"sellerId": "1"
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	UpdateProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestUpdateProductInvalidUserIdHeader
func TestUpdateProductInvalidUserIdHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/product", strings.NewReader(`{
		"amountAvailable": 22,				
		"cost": 150,			
		"productName": "Double Expresso",
		"sellerId": "1"
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	UpdateProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestUpdateProductDataAccessError
func TestUpdateProductDataAccessError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/product", strings.NewReader(`{
		"amountAvailable": 10,				
		"cost": 150,			
		"productName": "Double Expresso",
		"sellerId": "1"
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	testUpdProductDataError := testProduct{
		name:   "Test Update Product Data Error",
		userID: "seller-user-test",
		testID: "TestUpdateProductDataAccessError",
		mockUpdateThisProduct: func() {
			UpdateThisProduct = func(product *domains.Product) (*domains.Product, *config.HttpError) {
				return nil, config.DataAccessLayerError("")
			}
		},
	}
	testUpdProductDataError.mockUpdateThisProduct()
	UpdateProduct(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//TestUpdateProductValid
func TestUpdateProductValid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/product", strings.NewReader(`{
		"id": 1,
		"amountAvailable": 10,				
		"cost": 150,			
		"productName": "Double Expresso",
		"sellerId": "1"
	}`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	testUpdProductDataError := testProduct{
		name:   "Test Update Product Data Error",
		userID: "seller-user-test",
		testID: "TestUpdateProductDataAccessError",
		mockUpdateThisProduct: func() {
			UpdateThisProduct = func(product *domains.Product) (*domains.Product, *config.HttpError) {
				return &domains.Product{
					ID:              1,
					AmountAvailable: 10,
					Cost:            150,
					ProductName:     "Double Expresso",
					SellerId:        "1",
				}, nil
			}
		},
	}
	testUpdProductDataError.mockUpdateThisProduct()
	UpdateProduct(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

//TestDeleteProductInvalidProductId
func TestDeleteProductInvalidProductId(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/product/s", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "s"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	DeleteProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestDeleteProductNoID
func TestDeleteProductNoID(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/product", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	DeleteProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

//TestDeleteProductDataAccessError
func TestDeleteProductDataAccessError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/product/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	testDelProductDataError := testProduct{
		name:   "Test Update Product Data Error",
		userID: "seller-user-test",
		testID: "TestUpdateProductDataAccessError",
		mockDeleteThisProduct: func() {
			DeleteThisProduct = func(id int) *config.HttpError {
				return config.DataAccessLayerError("")
			}
		},
	}
	testDelProductDataError.mockDeleteThisProduct()
	DeleteProduct(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

//TestDeleteProductValid
func TestDeleteProductValid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/product/1", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("user_id", "seller-user-test")
	testDelProductDataError := testProduct{
		name:   "Test Update Product Data Error",
		userID: "seller-user-test",
		testID: "TestUpdateProductDataAccessError",
		mockDeleteThisProduct: func() {
			DeleteThisProduct = func(id int) *config.HttpError {
				return nil
			}
		},
	}
	testDelProductDataError.mockDeleteThisProduct()
	DeleteProduct(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
