package handler

import (
	"net/http"
	"strconv"

	"github.com/achsanit/go-assignment2/main/model"
	"github.com/achsanit/go-assignment2/main/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	GetOrders(c *gin.Context)
	CreateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
	UpdateOrder(c *gin.Context)
}

type orderHandlerImpl struct {
	service service.OrderService
}

func (orderHandler *orderHandlerImpl) GetOrders(c *gin.Context) {
	orders, err := orderHandler.service.GetOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (o *orderHandlerImpl) CreateOrder(ctx *gin.Context) {
	newOrder := model.OrderItems{}
	if err := ctx.Bind(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res, err := o.service.CreateOrder(ctx, newOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "data": res})
}

func (o *orderHandlerImpl) UpdateOrder(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newOrder := model.OrderItems{}
	if err := ctx.Bind(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res, err := o.service.UpdateOrder(ctx, uint64(orderId), newOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "data": res})
}

func (o *orderHandlerImpl) DeleteOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := o.service.DeleteOrder(c, uint64(orderId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func NewOrderHandler(svc service.OrderService) OrderHandler {
	return &orderHandlerImpl{
		service: svc,
	}
}
