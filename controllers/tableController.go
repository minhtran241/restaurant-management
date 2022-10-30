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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred while listing tables"},
			)
			return
		}
		var allTables []bson.M

		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		tableId := c.Param("table_id")
		var table models.Table
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "table was not found"})
			return
		} else if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred when fetching the table"},
			)
			return
		}
		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, insertErr := tableCollection.InsertOne(ctx, table)
		if insertErr != nil {
			msg := fmt.Sprintf("Failed to created table")
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var table models.Table

		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.Number_of_guests != nil {
			updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})
		}

		if table.Table_number != nil {
			updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", table.Updated_at})

		upsert := true
		filter := bson.M{"table_id": tableId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := tableCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)

		if err != nil {
			msg := "Failed to update table"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
