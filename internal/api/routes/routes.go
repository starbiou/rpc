package routes

import (
	"billing/internal/db/ent"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"os"
)

// RouteConfig holds the configuration for route setup
type RouteConfig struct {
	Router    *gin.Engine
	Client    *ent.Client
	Validator *validator.Validate
}

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
	router.Use(gin.LoggerWithWriter(os.Stdout)) // logging middleware
	router.Use(gin.Recovery())                  // Recovery middleware for panic handling
	validate := validator.New()

	config := RouteConfig{
		Router:    router,
		Client:    client,
		Validator: validate,
	}

	SetupReceiptRoutes(config)

	return router
}
