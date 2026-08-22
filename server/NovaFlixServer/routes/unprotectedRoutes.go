package routes

import (
	controller "github.com/Amaan0907/NovaFlix/server/NovaFlixServer/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine, client *mongo.Client) {

	router.GET("/movies", controller.GetMovies(client ))
	router.POST("/registeruser", controller.RegisterUser(client))
	router.POST("/login", controller.LoginUser(client))
	router.GET("/genres", controller.GetGenres(client ))

}
