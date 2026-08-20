package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/database"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/models"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenConnection("users")

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Check for an existing user before doing the (relatively expensive) hash.
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash the password"})
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			// If email has a unique index, a duplicate-key race lands here too.
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}


func LoginUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		var userLogin models.UserLogin

		if err:=c.ShouldBindJSON(&userLogin);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid Input data"})
			return
		}

		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()
		var foundUser models.User

		err:=userCollection.FindOne(ctx,bson.M{"email":userLogin.Email}).Decode(&foundUser)

		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid user or password"})
			return 
		}

		err=bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Invalid email or password"})
			return 
		}
		
		token,refreshToken,err:=utils.GenerateAllTokens(foundUser.FirstName,foundUser.LastName,foundUser.UserID,foundUser.Email,foundUser.Role)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to generate Tokens"})
			return 
		}
		err=utils.UpdateAllTokens(foundUser.UserID,token,refreshToken)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update User"})
			return 
		}
		c.JSON(http.StatusOK,models.UserResponse{
			FirstName: foundUser.FirstName,
			LastName: foundUser.LastName,
			UserId: foundUser.UserID,
			Email: foundUser.Email,
			Role: foundUser.Role,
			Token: token,
			RefreshToken: refreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})

	}
}