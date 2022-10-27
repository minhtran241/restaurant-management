package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func MenuRoutes(in *gin.Engine) {
	in.GET("/menus", controller.GetMenus())
	in.GET("/menus/:menu_id", controller.GetMenu())
	in.POST("/menus", controller.CreateMenu())
	in.PATCH("/menus/:menu_id", controller.UpdateMenu())
}