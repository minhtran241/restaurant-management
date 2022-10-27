package routes

import (
	"github.com/gin-gonic/gin"

	controller "github.com/minhtran241/restaurant-management/controllers"
)

func UserRoutes(in *gin.Engine) {
	in.GET("/users", controller.GetUsers())
	in.GET("/users/:user_id", controller.GetUser())
	in.POST("/users/signup", controller.SignUp())
	in.POST("users/login", controller.Login())
}