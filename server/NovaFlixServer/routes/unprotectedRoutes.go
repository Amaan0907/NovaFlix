

package routes

import (
	controller "github.com/Amaan0907/NovaFlix/server/NovaFlixServer/controllers"
	
	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine){
	
	router.GET("/movies",controller.GetMovies())
	router.POST("/registeruser",controller.RegisterUser())
	router.POST("/login",controller.LoginUser())

}