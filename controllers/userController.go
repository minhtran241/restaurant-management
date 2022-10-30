package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/minhtran241/restaurant-management/database"
	"github.com/minhtran241/restaurant-management/helpers"
	"github.com/minhtran241/restaurant-management/models"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// GetUsers responds with the list of all users as JSON.
// GetUsers             godoc
//  @Summary      Get all users
//  @Description  Responds with the list of all users as JSON.
//  @Tags         users
//  @Produce      json
//  @Success      200  {array}  models.User
//  @Router       /users [get]
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

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
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

// GetUser responds with the user with provided ID as JSON.
// GetUser             godoc
//  @Summary      Get single user by ID
//  @Description  Responds with the user with provided ID as JSON.
//  @Tags         users
//  @Produce      json
//  @Success      200  {object}  models.User
//  @Router       /users/{user_id} [get]
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

// SignUp takes user's information, provides JWT and stores in DB.
// SignUp             godoc
//  @Summary      Create a new user.
//  @Description  Create a new user.
//  @Tags         users
//  @Produce      json
//  @Success      200  {object}  models.User
//  @Router       /users/signup [post]
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var user models.User
		// convert the JSON data coming from client to golang readable format
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// validate the data based on user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		// check if the email has already been used by another user
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists"})
			return
		}
		// hash password
		password := HashPassword(*user.Password)
		user.Password = &password
		// check if the phone no. has already been used by another user
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for phone number"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "phone number already exists"})
			return
		}
		// get extra details for user object - created_at, updated_at, ID
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		// generate token and refresh token (helpers package)
		token, refreshToken, _ := helpers.GenerateAllTokens(
			*user.Email, *user.First_name, *user.Last_name, user.User_id,
		)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		// insert new user into the database
		resultInsertNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("Failed to create user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// return status OK and result
		c.JSON(http.StatusOK, resultInsertNumber)
	}
}

// SignUp takes user's information, compares with information in DB and provides new JWT.
// SignUp             godoc
//  @Summary      Log a user in..
//  @Description  Log a user in.
//  @Tags         users
//  @Produce      json
//  @Success      200  {object}  models.User
//  @Router       /users/login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		// convert the login data coming from client to golang readable format
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// find a user with that email, see if user even exists
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found, login failed"})
			return
		} else if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "user not found, login failed"},
			)
			return
		}
		// verify password
		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// generate tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(
			*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id,
		)
		// update tokens - token and refresh token
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		// return status OK and result
		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, providePassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providePassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("incorrect password, login failed")
		check = false
	}
	return check, msg
}
