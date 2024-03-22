package model

import (
	"time"
)

type Order struct {
	OrderId      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
}

type Item struct {
	ItemId      int    `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	OrderId     int    `json:"order_id"`
}

type OrderItems struct {
	OrderedAt    time.Time `json:"ordered_at"`
	CustomerName string    `json:"customer_name"`
	Items        []Item    `json:"items"`
}
