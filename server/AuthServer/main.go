// Entry point of the application where the server is initialized.
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/naman2607/netflixClone/config"
	"github.com/naman2607/netflixClone/database"
	authRoute "github.com/naman2607/netflixClone/routes"
)

func initRouter(r *gin.Engine) {
	authRoute.InitAuthRoute(r)
}

func main() {

	serverConfig := config.GetInstance()
	serverConfig.InitServerConfig()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := database.InitPostgresDB(serverConfig)
	if err != nil {
		log.Println("Failed to open database ", err)
	}
	initRouter(r)
	r.Run(":" + serverConfig.GetServerPort())
	defer database.GetDB().Close()
}
