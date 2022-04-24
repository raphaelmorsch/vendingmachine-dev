package repositories

import (
	"errors"
	"vendingmachine/config"
	"vendingmachine/domains"

	"gorm.io/gorm"
)

func FinishPurchase(product domains.Product, userId string) error {
	var result = config.Database.Transaction(func(tx *gorm.DB) error {

		_, errUpdProduct := UpdateProduct(&product)
		//Reset User's Deposits
		errDelUD := DeleteUserDeposit(userId)

		if errUpdProduct != nil {
			return errors.New(errUpdProduct.Error)
		} else {
			if errDelUD != nil {
				return errors.New(errDelUD.Error)
			}
		}
		return nil
	})
	return result
}
