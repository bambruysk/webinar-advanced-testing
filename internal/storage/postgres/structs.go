package postgres

import (
	"time"

	"webinar-testing/pkg/models"
)

type ListResult struct {
	ID        uint64
	UserID    string
	Good      string
	Quantity  int32
	CreatedAt time.Time
}

func listResultsToOrder(results []ListResult, userID models.UserID) (order models.Order) {
	if len(results) == 0 {
		order.UserID = userID
		return order
	}

	order.Goods = make(map[models.GoodID]int)
	order.UserID = userID

	for _, res := range results {
		good := models.GoodID(res.Good)
		if qty, exist := order.Goods[good]; exist {
			order.Goods[good] = qty + int(res.Quantity)

			continue
		}

		order.Goods[good] = int(res.Quantity)
	}

	return order
}
