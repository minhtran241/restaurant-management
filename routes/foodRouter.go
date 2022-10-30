package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func FoodRoutes(in *gin.Engine) {
	in.GET("/foods", controller.GetFoods())
	in.GET("/foods/:food_id", controller.GetFood())
	in.POST("/foods", controller.CreateFood())
	in.PATCH("/foods/:food_id", controller.UpdateFood())
}
