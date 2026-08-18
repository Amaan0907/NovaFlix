package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	controller "github.com/Amaan0907/NovaFlix/server/NovaFlixServer/controllers"
)

func main() {
	
	router:=gin.Default()

	router.GET("/hello",func(c *gin.Context) {
		c.String(http.StatusOK,"Hello, NovaFlix")
	})
	router.GET("/movies",controller.GetMovies())
	
	if err:=router.Run(":3000");err!=nil{
		fmt.Println("Failed to start server",err)
		panic(err)
	}


}