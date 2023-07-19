package repository

import (
	"context"
	"mediasoft-statistics/internal/model"
	"time"
)

type OrderRepository interface {
	Create(model *model.Order) error
	ListOrderItems(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.OrderProduct, error)
}
