package model

import (
	"github.com/google/uuid"
	"time"
)

type OrderProduct struct {
	CreatedAt   time.Time `db:"created_at"`
	Count       int64     `db:"count"`
	ProductUuid uuid.UUID `db:"product_uuid"`
}
