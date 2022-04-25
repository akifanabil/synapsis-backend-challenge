package interfaces

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name     		string
	Category		string
	Description     string
	Amount	 		int
	Price			int
	Cart  			[]Cart `gorm:"foreignKey:ProductID"`
	Transactions	[]Transaction `gorm:"foreignKey:CustomerID"`
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