package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/minhtran241/restaurant-management/database"
	"github.com/minhtran241/restaurant-management/models"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing users"})
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		userId := c.Param("user_id")
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "user was not found"})
			return
		} else if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "error occurred when fetching user"},
			)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert the JSON data coming from client to golang readable format

		// validate the data based on user struct

		// check if the email has already been used by another user

		// hash password

		// check if the phone no. has already been used by another user

		// get extra details for user object - created_at, updated_at, ID

		// generate token and refresh token (helpers package)

		// insert new user into the database

		// return status OK and result

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// convert the login data coming from client to golang readable format

		// find a user with that email, see if user even exists

		// verify password

		// generate tokens

		// update tokens - token and refresh token

		// return status OK and result
	}
}

func HashPassword(password string) string {

}

func VerifyPassword(password, providePassword string) (bool, string) {

}
