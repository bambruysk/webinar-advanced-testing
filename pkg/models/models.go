package models

type UserID string

type GoodID string

type Order struct {
	UserID UserID
	// Goods товары в заказе и их количество.
	Goods map[GoodID]int
}

func (o Order) Validate() interface{} {
	return nil
}
