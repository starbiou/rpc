package routes

import (
	"github.com/gin-gonic/gin"
	"receipt-processor/ent"
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

	//router.GET("/", controllers.IndexView)
	//router.GET("/items", func(context *gin.Context) { controllers.TodoItems(context, client) })
	//router.POST("/item", func(context *gin.Context) { controllers.CreateTodoItems(context, client) })
	//router.PUT("/item/:id", func(context *gin.Context) { controllers.UpdateTodoItem(context, client) })
	//router.DELETE("/item/:id", func(context *gin.Context) { controllers.DeleteTodoItem(context, client) })
	//router.PUT("/items", func(context *gin.Context) { controllers.UpdateTodoItems(context, client) })
	//router.DELETE("/items", func(context *gin.Context) { controllers.DeleteTodoItems(context, client) })
	//router.GET("/items/filter", func(context *gin.Context) { controllers.FilterTodoItems(context, client) })

	return router
}
