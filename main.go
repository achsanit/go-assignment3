package main

import (
	"github.com/achsanit/go-assignment2/main/handler"
	"github.com/achsanit/go-assignment2/main/infrastructure"
	"github.com/achsanit/go-assignment2/main/repository"
	"github.com/achsanit/go-assignment2/main/router"
	"github.com/achsanit/go-assignment2/main/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	postgres := infrastructure.NewSqlPostgres()

	v1 := r.Group("/v1")
	{
		orders := v1.Group("/orders")
		{
			orderRepo := repository.NewOrdeQuery(postgres)
			orderService := service.NewOrderService(orderRepo)
			orderHandler := handler.NewOrderHandler(orderService)
			orderRouter := router.NewOrderRouter(orders, orderHandler)

			orderRouter.Mount()
		}
	}

	r.Run(":8001")
}
