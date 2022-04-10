package repositories

import (
	"errors"
	"vendingmachine/src/config"
	"vendingmachine/src/domains"

	"gorm.io/gorm"
)

func FindUserDeposit(userId string) (*domains.UserDeposit, *config.HttpError) {

	var userDeposit domains.UserDeposit
	err := config.Database.First(&userDeposit, "user_id", userId).Error

	if err != nil {
		return nil, config.DataAccessLayerError(err.Error())
	}

	return &userDeposit, nil

}
func MakeUserDeposit(userDeposit *domains.UserDeposit) (*domains.UserDeposit, *config.HttpError) {

	var updDeposit domains.UserDeposit
	err := config.Database.First(&updDeposit, "user_id", userDeposit.UserId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userDeposit.Balance = userDeposit.Deposit
			e := config.Database.Create(&userDeposit)
			if e.Error != nil {
				return nil, config.DataAccessLayerError(e.Error.Error())
			}
			return userDeposit, nil
		}

	} else {
		updDeposit.Balance = updDeposit.Balance + userDeposit.Deposit
		e := config.Database.Save(&updDeposit)

		if e.Error != nil {
			return nil, config.DataAccessLayerError(e.Error.Error())
		}

		return &updDeposit, nil

	}

	return nil, nil

}

func DeleteUserDeposit(id string) *config.HttpError {

	var userDeposit domains.UserDeposit
	err := config.Database.First(&userDeposit, "user_id", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return config.NotFoundError()
		} else {
			return config.DataAccessLayerError(err.Error())
		}
	}
	e := config.Database.Delete(&userDeposit, "user_id", id)

	if e.Error != nil {
		return config.DataAccessLayerError(e.Error.Error())
	}

	return nil

}
