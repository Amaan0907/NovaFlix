package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	
	router:=gin.Default()

	router.GET("/hello",func(ctx *gin.Context) {
		ctx.String(http.StatusOK,"Hello, NovaFlix")
	})

	if err:=router.Run(":3000");err!=nil{
		fmt.Println("Failed to start server",err)
		panic(err)
	}


}