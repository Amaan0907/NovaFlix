package routes

import (
	controller "github.com/Amaan0907/NovaFlix/server/NovaFlixServer/controllers"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine,client *mongo.Client){
	router.Use(middleware.AuthMiddleware())

	router.GET("/movie/:imdb_id",controller.GetMovie(client))
	router.POST("/addmovie",controller.AddMovie(client))
	router.GET("/recommendedmovies",controller.GetRecommendedMovies(client))
	router.PATCH("/updatereview/:imdb_id",controller.AdminReviewUpdate(client))

}