package interfaces

import "gorm.io/gorm"

type Chart struct {
	gorm.Model
	ProductID     uint
	CustomerID    uint
	Amount		  int
}

type ChartResponse struct {
	ChartID			uint
	ProductID		uint
	Name			string
	Price			int
	Category		string
	Description     string
	Amount	 		int
}

type Charts struct {
	Charts []ChartResponse
}