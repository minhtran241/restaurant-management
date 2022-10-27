package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func OrderItemRoutes(in *gin.Engine) {
	in.GET("/orderItems", controller.GetOrderItems())
	in.GET("/orderItems/:orderItem_id", controller.GetOrderItem())
	in.GET("/orderItems-order/:order_id", controller.GetOrderItemsByOrder())
	in.POST("/orderItems", controller.CreateOrderItem())
	in.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItem())
}