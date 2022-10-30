package controllers

import (
	"context"
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

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := orderItemCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred while listing ordered items"},
			)
			return
		}
		var allOrderItems []bson.M

		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem
		err := orderItemCollection.FindOne(
			ctx,
			bson.M{"orderItem_id": orderItemId},
		).Decode(&orderItem)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "ordered item was not found"})
			return
		} else if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred when fetching the ordered item"},
			)
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred while listing ordered items by order ID"},
			)
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(
			time.RFC3339, time.Now().Format(time.RFC3339),
		)

		orderItemsToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, insertOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")
		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", orderItem.Unit_price})
		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true
		filter := bson.M{"order_item_id": orderItemId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)

		if err != nil {
			msg := "Failed to update the ordered item"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
