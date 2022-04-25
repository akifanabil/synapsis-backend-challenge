package interfaces

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	CustomerID  uint
	ProductID   uint
	Amount		int
	Totals   	int
}