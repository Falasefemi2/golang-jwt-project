// package main

// import (
// 	"os"

// 	"github.com/falasefemi2/golang-jwt-project/routes"
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8000"
// 	}

// 	router := gin.New()
// 	router.Use(gin.Logger())

// 	routes.AuthRoutes(router)
// 	routes.UserRoutes(router)

// 	router.GET("/api-1", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"success": "Access granted for api-1"})
// 	})

// 	router.GET("/api-2", func(c *gin.Context) {
// 		c.JSON(200, gin.H{"success": "Access granted for api-2"})
// 	})

// 	// Start the server
// 	router.Run(":" + port)
// }

package main

import (
	"os"

	"github.com/falasefemi2/golang-jwt-project/middleware"
	"github.com/falasefemi2/golang-jwt-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	// Public routes
	routes.AuthRoutes(router)

	// Protected routes
	router.Use(middleware.Authenticate())
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	// Start the server
	router.Run(":" + port)
}
