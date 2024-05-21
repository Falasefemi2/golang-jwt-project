package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/falasefemi2/golang-jwt-project/database"
	"github.com/falasefemi2/golang-jwt-project/helpers"
	"github.com/falasefemi2/golang-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// func HashPassword(password string) string {
// 	 bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return string(bytes)
// }

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprint("password incorrect")
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Ensure this is called only once to release resources

		var user models.User

		// Bind JSON to user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the user struct
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if email exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking email existence"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}

		// Check if phone exists
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking phone existence"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
			return
		}

		// Hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		// Set user metadata
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Generate tokens (replace with your actual token generation function)
		token, refreshToken, err := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
			return
		}
		user.Token = &token
		user.Refresh_token = &refreshToken

		// Insert user into collection
		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}

		// Respond with success
		c.JSON(http.StatusOK, gin.H{"message": "User successfully registered"})
	}
}

// func Login() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		var user models.User
// 		var foundUser models.User

// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
// 			return
// 		}

// 		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
// 		defer cancel()
// 		if passwordIsValid != true {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}

// 		if foundUser.Email == nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
// 		}
// 		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
// 		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
// 		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, foundUser)
// 	}
// }

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Move the defer cancel() statement to the top

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		// Generate new tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)

		// Update tokens in the database
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// Update the foundUser struct with the new tokens
		foundUser.Token = &token
		foundUser.Refresh_token = &refreshToken

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helpers.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

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

		matchStage := bson.D{{"$match", bson.D{}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}},
		}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
