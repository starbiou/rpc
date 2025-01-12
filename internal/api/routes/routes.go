package routes

import (
	"billing/internal/db/ent"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(client *ent.Client) *gin.Engine {
	router := gin.Default()
	// CORS configuration
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	SetupReceiptRoutes(router, client)

	return router
}
