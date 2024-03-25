package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/achsanit/go-assignment2/main/infrastructure"
	"github.com/achsanit/go-assignment2/main/model"
)

type OrderQuery interface {
	GetOrders(c context.Context) ([]model.Order, error)
	CreateOrder(c context.Context, order model.OrderItems) (model.OrderItems, error)
	DeleteOrder(c context.Context, orderId uint64) error
	UpdateOrder(c context.Context, orderId uint64, order model.OrderItems) (model.OrderItems, error)
}

type orderQueryImpl struct {
	db infrastructure.SqlPostgres
}

func NewOrdeQuery(db infrastructure.SqlPostgres) OrderQuery {
	return &orderQueryImpl{
		db: db,
	}
}

func (q *orderQueryImpl) CreateOrder(c context.Context, order model.OrderItems) (newOrder model.OrderItems, errMessage error) {
	db := q.db.GetConnection()

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return model.OrderItems{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	o := model.Order{
		CustomerName: order.CustomerName,
		OrderedAt:    time.Now(),
	}
	var orderId int
	if err := db.QueryRow(
		"INSERT INTO orders (customer_name, ordered_at) VALUES ($1, $2) RETURNING order_id",
		o.CustomerName,
		o.OrderedAt,
	).Scan(&orderId); err != nil {
		tx.Rollback()
		panic(err)
	}

	resItem := []model.Item{}
	for _, data := range order.Items {
		data.OrderId = orderId

		newItem := model.Item{}
		log.Println(newItem)
		err := db.
			QueryRow(
				"INSERT INTO items (item_code, description, quantity, order_id) VALUES ($1, $2, $3, $4) RETURNING * ",
				data.ItemCode,
				data.Description,
				data.Quantity,
				data.OrderId).
			Scan(
				&newItem.ItemId,
				&newItem.ItemCode,
				&newItem.Description,
				&newItem.Quantity,
				&newItem.OrderId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		resItem = append(resItem, newItem)
	}

	return model.OrderItems{
		CustomerName: order.CustomerName,
		Items:        resItem,
	}, nil
}

func (q *orderQueryImpl) GetOrders(c context.Context) ([]model.Order, error) {
	db := q.db.GetConnection()
	if db == nil {
		return []model.Order{}, errors.New("db nil")
	}

	orders := []model.Order{}
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt); err != nil {
			return []model.Order{}, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (q *orderQueryImpl) DeleteOrder(c context.Context, orderId uint64) error {
	db := q.db.GetConnection()

	//find order
	orders := []model.Order{}
	rows, err := db.Query("SELECT * FROM orders WHERE order_id = $1", orderId)
	if err != nil {
		return err
	}

	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt); err != nil {
			return err
		}
		orders = append(orders, order)
	}

	if len(orders) < 1 {
		return errors.New("data order not found")
	}

	// delete order
	res, err := db.Exec("DELETE FROM orders WHERE order_id = $1", orderId)
	if err != nil {
		return err
	}

	val, _ := res.RowsAffected()
	if val <= 0 {
		return errors.New("no rows deleted")
	}

	return nil
}

func (q *orderQueryImpl) UpdateOrder(c context.Context, orderId uint64, order model.OrderItems) (newOrder model.OrderItems, errMessage error) {
	db := q.db.GetConnection()

	// check get order data
	orders := []model.Order{}
	rows, err := db.Query("SELECT * FROM orders WHERE order_id = $1", orderId)
	if err != nil {
		return model.OrderItems{}, err
	}

	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt); err != nil {
			return model.OrderItems{}, err
		}
		orders = append(orders, order)
	}

	if len(orders) < 1 {
		return model.OrderItems{}, errors.New("data order not found")
	}

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return model.OrderItems{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	o := model.Order{
		CustomerName: order.CustomerName,
		OrderedAt:    time.Now(),
	}
	if _, err := db.Exec(
		"UPDATE orders SET customer_name = $1, ordered_at = $2 WHERE order_id = $3",
		o.CustomerName,
		o.OrderedAt,
		orderId,
	); err != nil {
		log.Default().Println(order, "::", orderId)
		tx.Rollback()
		panic(err)
	}

	// update each item
	log.Println(order.Items)
	for _, data := range order.Items {
		// check get item
		log.Println(data)
		rows, err := db.Query("SELECT * FROM items WHERE item_id = $1", data.ItemId)
		if err != nil {
			log.Println("select item error::", data.ItemId, err.Error())
			continue
		}

		items := []model.Item{}
		for rows.Next() {
			var item model.Item
			if err := rows.Scan(&item.ItemId, &item.ItemCode, &item.Description, &item.Quantity, &item.OrderId); err != nil {
				return model.OrderItems{}, err
			}
			items = append(items, item)
		}

		if len(items) < 1 {
			log.Println("item not found::", data.ItemId)
			continue
		}

		// update item
		if _, err := db.Exec(
			"UPDATE items SET item_code = $1, description = $2, quantity = $3 WHERE item_id = $4",
			data.ItemCode,
			data.Description,
			data.Quantity,
			data.ItemId,
		); err != nil {
			log.Println("select item error::", err.Error())
			continue
		}
	}

	return model.OrderItems{
		CustomerName: order.CustomerName,
		Items:        order.Items,
	}, nil
}
