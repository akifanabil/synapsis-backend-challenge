package interfaces

type Cart struct {
	ProductID     uint `gorm:"primaryKey"`
	CustomerID    uint `gorm:"primaryKey"`
	Amount		  int
}

type CartResponse struct {
	ProductID		uint
	Name			string
	Price			int
	Category		string
	Description     string
	Amount	 		int
}

type Carts struct {
	Carts []CartResponse
}