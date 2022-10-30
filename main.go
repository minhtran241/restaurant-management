package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/minhtran241/restaurant-management/database"
	_ "github.com/minhtran241/restaurant-management/docs"
	"github.com/minhtran241/restaurant-management/middleware"
	"github.com/minhtran241/restaurant-management/routes"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// @title Restaurant Management Service
// @version 1.0
// @description Restaurant Management Service API using Gin framework.
// @termsOfService http://swagger.io/terms/

// @contact.name Minh Tran
// @contact.email trqminh24@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	// router.Use(middleware.CORSMiddleware())
	router.Use(cors.Default())
	router.Use(gin.Logger())

	routes.UserRoutes(router)

	router.GET("/", HealthCheck)
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)
}

// HealthCheck godoc
//  @Summary      Show the status of server.
//  @Description  get the status of server.
//  @Tags         root
//  @Accept       */*
//  @Produce      json
//  @Success      200  {object}  map[string]interface{}
//  @Router       / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}
	c.JSON(http.StatusOK, res)
}
