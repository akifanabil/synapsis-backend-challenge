package interfaces

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     		string
	Category		string
	Description     string
	Amount	 		int
	Price			int
	Chart  			[]Chart `gorm:"foreignKey:ProductID"`
}

type ProductResponse struct {
	ID				uint
	Name     		string
	Category		string
	Description     string
	Amount	 		int
	Price			int
}

type Products struct {
	Products []ProductResponse
}