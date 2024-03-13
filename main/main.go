package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"rhzx3519/unique-key-generator/pool"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	r := gin.Default()

	p := pool.NewPool()
	p.Run(context.TODO())

	v1 := r.Group("/v1")
	{
		v1.GET("/key", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"key": p.Key(),
			})
		})
		v1.GET("/existed", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"existed": p.Existed(ctx.Query("key")),
			})
		})

	}

	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	r.Run(port)
}
