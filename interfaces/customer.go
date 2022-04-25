package interfaces

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name     		string
	Email    		string
	Password 		string
	Gender   		string
	Balance  		int
	Cart  			[]Cart `gorm:"foreignKey:CustomerID"`
	Transactions	[]Transaction `gorm:"foreignKey:CustomerID"`
}

type CustomerResponse struct {
	ID     	uint
	Name   	string
	Email  	string
	Gender 	string
	Balance int
}

type AuthResponse struct {
	Data 			CustomerResponse
	ExpiresIn_hour	uint
	Jwt				string
	Message			string
}