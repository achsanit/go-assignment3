package service

import (
	"context"

	"github.com/achsanit/go-assignment2/main/model"
	"github.com/achsanit/go-assignment2/main/repository"
)

type OrderService interface {
	GetOrders(c context.Context) ([]model.Order, error)
	CreateOrder(c context.Context, order model.OrderItems) (model.OrderItems, error)
	DeleteOrder(c context.Context, orderId uint64) error
	UpdateOrder(c context.Context, orderId uint64, order model.OrderItems) (model.OrderItems, error)
}

type orderServiceImpl struct {
	repository repository.OrderQuery
}

func NewOrderService(repo repository.OrderQuery) OrderService {
	return &orderServiceImpl{
		repository: repo,
	}
}

func (order *orderServiceImpl) GetOrders(ctx context.Context) ([]model.Order, error) {
	orders, err := order.repository.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderServiceImpl) CreateOrder(ctx context.Context, order model.OrderItems) (model.OrderItems, error) {
	newOrder, err := o.repository.CreateOrder(ctx, order)
	if err != nil {
		return model.OrderItems{}, err
	}

	return newOrder, nil
}

func (o *orderServiceImpl) DeleteOrder(ctx context.Context, orderId uint64) error {
	if err := o.repository.DeleteOrder(ctx, orderId); err != nil {
		return err
	}

	return nil
}

func (o *orderServiceImpl) UpdateOrder(ctx context.Context, orderId uint64, order model.OrderItems) (model.OrderItems, error) {
	res, err := o.repository.UpdateOrder(ctx, orderId, order)
	if err != nil {
		return model.OrderItems{}, err
	}

	return res, nil
}
