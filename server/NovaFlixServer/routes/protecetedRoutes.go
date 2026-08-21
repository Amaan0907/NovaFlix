package routes

import (
	controller "github.com/Amaan0907/NovaFlix/server/NovaFlixServer/controllers"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine){
	router.Use(middleware.AuthMiddleware())

	router.GET("/movie/:imdb_id",controller.GetMovie())
	router.POST("/addmovie",controller.AddMovie())
	router.GET("/recommendedmovies",controller.GetRecommendedMovies())
	router.PATCH("/updatereview/:imdb_id",controller.AdminReviewUpdate())

}