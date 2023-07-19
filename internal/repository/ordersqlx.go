package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"mediasoft-statistics/internal/model"
	"time"
)

type OrderSqlx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *OrderSqlx {
	return &OrderSqlx{db: db}
}

func (r *OrderSqlx) Create(order *model.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	const q = `insert into orders (uuid, user_uuid) values(:uuid, :user_uuid)`
	_, err = tx.NamedExec(q, order)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	var orderItemq = `insert into order_item(order_uuid, count, product_uuid) values ($1,$2,$3)`
	for _, orderItems := range [][]*model.OrderItem{order.Salads, order.Soups, order.Drinks, order.Desserts, order.Meats, order.Garnishes} {
		for _, orderItem := range orderItems {
			_, err = tx.Exec(orderItemq, order.Uuid, orderItem.Count, orderItem.ProductUuid)
			if err != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					return errors.Join(err, errRollback)
				}
			}
		}
	}

	return tx.Commit()
}

func (r *OrderSqlx) ListOrderItems(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.OrderProduct, error) {
	const q = `SELECT o.created_at, oi.count, oi.product_uuid 
				FROM order_item oi
    			JOIN orders o ON o.uuid = oi.order_uuid
				WHERE date_trunc('day',o.created_at) >= date_trunc('day', $1::timestamp)
				AND date_trunc('day',o.created_at) <= date_trunc('day', $2::timestamp);`
	var list []*model.OrderProduct
	err := r.db.SelectContext(ctx, &list, q, startDate, endDate)
	return list, err
}
