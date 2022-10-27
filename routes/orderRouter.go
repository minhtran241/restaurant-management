package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func OrderRoutes(in *gin.Engine) {
	in.GET("/orders", controller.GetOrders())
	in.GET("/orders/:order_id", controller.GetOrder())
	in.POST("/orders", controller.CreateOrder())
	in.PATCH("/orders/:order_id", controller.UpdateOrder())
}