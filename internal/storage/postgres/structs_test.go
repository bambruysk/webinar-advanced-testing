package postgres

import (
	"reflect"
	"testing"
	"time"

	"webinar-testing/pkg/models"
)

func Test_listResultsToOrder(t *testing.T) {
	type args struct {
		results []ListResult
		userID  models.UserID
	}
	tests := []struct {
		name      string
		args      args
		wantOrder models.Order
	}{
		{
			name: "Test 1",
			args: args{
				userID: models.UserID("user"),
				results: []ListResult{
					{
						ID:        1,
						UserID:    "user",
						Good:      "g1",
						Quantity:  12,
						CreatedAt: time.Time{},
					},
					{
						ID:        2,
						UserID:    "user",
						Good:      "g1",
						Quantity:  -5,
						CreatedAt: time.Time{},
					},
				},
			},
			wantOrder: models.Order{
				UserID: models.UserID("user"),
				Goods: map[models.GoodID]int{
					"g1": 7,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOrder := listResultsToOrder(tt.args.results, tt.args.userID); !reflect.DeepEqual(gotOrder, tt.wantOrder) {
				t.Errorf("listResultsToOrder() = %v, want %v", gotOrder, tt.wantOrder)
			}
		})
	}
}
