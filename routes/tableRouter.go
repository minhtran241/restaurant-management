package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func TableRoutes(in *gin.Engine) {
	in.GET("/tables", controller.GetTables())
	in.GET("/tables/:table_id", controller.GetTable())
	in.POST("/tables", controller.CreateTable())
	in.PATCH("/tables/:table_id", controller.UpdateTable())
}