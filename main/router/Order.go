package router

import (
	"github.com/achsanit/go-assignment2/main/handler"
	"github.com/gin-gonic/gin"
)

type OrderRouter interface {
	Mount()
}

type orderRouterImpl struct {
	version *gin.RouterGroup
	handler handler.OrderHandler
}

func NewOrderRouter(v *gin.RouterGroup, h handler.OrderHandler) OrderRouter {
	return &orderRouterImpl{version: v, handler: h}
}

func (order *orderRouterImpl) Mount() {
	order.version.GET("", order.handler.GetOrders)
	order.version.POST("", order.handler.CreateOrder)
	order.version.DELETE("/:order_id", order.handler.DeleteOrder)
	order.version.PUT("/:order_id", order.handler.UpdateOrder)
}
