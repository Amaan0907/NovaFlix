package main

import (
	"fmt"

	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/routes"
	
	"github.com/gin-gonic/gin"
)

func main() {
	
	router:=gin.Default()
	routes.SetupUnprotectedRoutes(router)
	routes.SetupProtectedRoutes(router)
	
	if err:=router.Run(":3000");err!=nil{
		fmt.Println("Failed to start server",err)
		panic(err)
	}


}