package main

import (
	"backend/middleware"
	routes "backend/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
	routes.UserRoutes(router)
	routes.PostRoutes(router)
	routes.LikeRoutes(router)

	router.Use(middleware.Authentication())
	// API-2
	router.GET("/api-1", func(c *gin.Context) {

		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
