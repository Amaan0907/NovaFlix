package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/database"
	"github.com/Amaan0907/NovaFlix/server/NovaFlixServer/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	var client *mongo.Client = database.Connect()


	if err:=client.Ping(context.Background(),nil);err!=nil{
		log.Fatalf("Failed to reach server: %v",err)
	}
	defer func(){
		err:=client.Disconnect(context.Background())
		if err!=nil{
			log.Fatalf("Failed to disconnect: %v",err)
		}
	}()

	
	routes.SetupUnprotectedRoutes(router,client)
	routes.SetupProtectedRoutes(router,client)

	if err := router.Run(":3000"); err != nil {
		fmt.Println("Failed to start server", err)
		panic(err)
	}

}
