package repositories

import (
	"errors"
	"vendingmachine/config"
	"vendingmachine/domains"

	"gorm.io/gorm"
)

func SaveProduct(product *domains.Product) (*domains.Product, *config.HttpError) {

	e := config.Database.Create(&product)

	if e.Error != nil {
		return nil, config.DataAccessLayerError(e.Error.Error())
	}

	return product, nil
}

func UpdateProduct(product *domains.Product) (*domains.Product, *config.HttpError) {

	var updProduct domains.Product
	config.Database.First(&updProduct, product.ID)

	updProduct.AmountAvailable = product.AmountAvailable
	updProduct.Cost = product.Cost
	updProduct.Cost = product.Cost
	updProduct.ProductName = product.ProductName
	updProduct.SellerId = product.SellerId

	e := config.Database.Save(&updProduct)

	if e.Error != nil {
		return nil, config.DataAccessLayerError(e.Error.Error())
	}

	return &updProduct, nil
}

func FindProductById(id int) *domains.Product {
	var product domains.Product

	config.Database.First(&product, id)

	return &product
}

func FindAllProducts() []domains.Product {

	var products []domains.Product

	config.Database.Find(&products)

	return products
}

func DeleteProduct(id int) *config.HttpError {

	var product domains.Product
	err := config.Database.First(&product, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return config.NotFoundError()
		} else {
			return config.DataAccessLayerError(err.Error())
		}
	}
	e := config.Database.Delete(&product, id)

	if e.Error != nil {
		return config.DataAccessLayerError(e.Error.Error())
	}

	return nil

}
