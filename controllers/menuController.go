package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/minhtran241/restaurant-management/database"
	"github.com/minhtran241/restaurant-management/models"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

// GetMenus responds with the list of all menus as JSON.
// GetMenus             godoc
//  @Summary      Get all menus
//  @Description  Responds with the list of all menus as JSON.
//  @Tags         menus
//  @Produce      json
//  @Success      200  {array}  models.Menu
//  @Router       /menus [get]
func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred while listing menu items"},
			)
			return
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

// GetMenu responds with the menu with provided ID as JSON.
// GetMenu             godoc
//  @Summary      Get single menu by ID
//  @Description  Responds with the menu with provided ID as JSON.
//  @Tags         menus
//  @Produce      json
//  @Success      200  {object}  models.Menu
//  @Router       /menus/{menu_id} [get]
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		menuId := c.Param("menu_id")
		var menu models.Menu
		err := foodCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu was not found"})
			return
		} else if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred when fetching the menu"},
			)
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

// CreateMenu takes a menu JSON and store in DB.
// CreateMenu             godoc
//  @Summary      Store a new menu
//  @Description  Takes a menu JSON and store in DB. Return saved JSON.
//  @Tags         menus
//  @Produce      json
//  @Success      200  {object}  models.Menu
//  @Router       /menus [post]
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("Failed to created menu item")
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// UpdateMenu takes a menu JSON and update menu stored in DB.
// UpdateMenu             godoc
//  @Summary      Update a menu
//  @Description  Takes a menu JSON and update menu stored in DB. Return saved JSON.
//  @Tags         menus
//  @Produce      json
//  @Success      200  {object}  models.Menu
//  @Router       /menus/{menu_id} [patch]
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}

		var updateObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil {
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "invalid time"
				c.JSON(http.StatusBadRequest, gin.H{"error": msg})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})

			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
			}
			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)

			if err != nil {
				msg := "Failed to update the menu"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			c.JSON(http.StatusOK, result)
		}

	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}
