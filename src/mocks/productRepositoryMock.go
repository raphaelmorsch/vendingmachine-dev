package mocks

import (
	"vendingmachine/src/config"
	"vendingmachine/src/domains"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindProductById(id int) *domains.Product {

	ret := m.Called(id)

	var r0 *domains.Product
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domains.Product)
	}

	return r0

}

func (m *MockProductRepository) UpdateProduct(product *domains.Product) (*domains.Product, *config.HttpError) {
	ret := m.Called(product)

	var r0 *domains.Product
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domains.Product)
	}

	var r1 *config.HttpError

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*config.HttpError)
	}

	return r0, r1
}

func (m *MockProductRepository) FindAllProducts() []domains.Product {
	ret := m.Called()

	var r0 []domains.Product
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]domains.Product)
	}
	return r0

}

func (m *MockProductRepository) DeleteProduct(id int) *config.HttpError {
	ret := m.Called(id)

	var r0 *config.HttpError
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*config.HttpError)
	}
	return r0

}
