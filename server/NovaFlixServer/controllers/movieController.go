package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/database"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenConnection("movies")

var validate=validator.New()
func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		var movies []models.Movie
		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch movies"})
				return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to decode movies"})
			return

		}
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}
		var movie models.Movie

		err:= movieCollection.FindOne(ctx,bson.M{"imdb_id":movieID}).Decode(&movie)
		if err!=nil{
			c.JSON(http.StatusNotFound,gin.H{"error":"Movie not Found"})
			return


		}

		c.JSON(http.StatusOK,movie)
	}
	
}

func AddMovie() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		var movie models.Movie

		if err:=c.ShouldBindJSON(&movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Movie Details are required to add Movie"})
			return
		}

		if err:=validate.Struct(movie);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"Validation Failed","details":err.Error()})
			return 
		}
		result,err:=movieCollection.InsertOne(ctx,movie)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to add Movie"})
			return
		}

		c.JSON(http.StatusCreated,result)

	}
}
