package main

import (
	"backend/middleware"
	routes "backend/router"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

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

	fmt.Println("Starting server on the port 8080")
	router.Run(":" + port)
	//log.Fatal(http.ListenAndServe(":8080", r))
}

//{
//"first_name" : "Zico",
//"last_name" : "Tjia",
//"username": "TjiaZico",
//"email": "zicotjia@gmail.com",
//"password": "france123"
//}
