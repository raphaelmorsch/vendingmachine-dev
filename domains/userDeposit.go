package domains

import "gorm.io/gorm"

type UserDeposit struct {
	gorm.Model

	Deposit int `json:"deposit"`

	Balance int `json:"balance"`

	UserId string `gorm:"primaryKey;autoIncrement:false" json:"userId"`
}
