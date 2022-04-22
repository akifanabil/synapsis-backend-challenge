package interfaces

import "gorm.io/gorm"

type Chart struct {
	gorm.Model
	ProductID     uint
	CustomerID    uint
	Amount		  int
}