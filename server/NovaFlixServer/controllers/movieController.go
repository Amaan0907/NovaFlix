package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/database"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection=database.OpenConnection("movies")
func GetMovies() gin.HandlerFunc{
	return func (c *gin.Context){
		ctx,cancel:=context.WithTimeout(context.Background(),100*time.Second)

		defer cancel()
		var movies []models.Movie
		cursor,err:= movieCollection.Find(ctx, bson.M{})

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Failed to fetch movies"})

		}
		defer cursor.Close(ctx)

		if err=cursor.All(ctx,&movies);err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"Error":"Failed to decode movies"})

		}
		c.JSON(http.StatusOK,movies)
	}
}















