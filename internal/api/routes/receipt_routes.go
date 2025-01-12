package routes

import (
	"billing/internal/api/controllers"
	"billing/internal/core/services"
	"billing/internal/db/ent"
	"billing/internal/db/repositories"
	"github.com/gin-gonic/gin"
)

// SetupReceiptRoutes defines the routes related to receipt processing.
func SetupReceiptRoutes(router *gin.Engine, client *ent.Client) {
	receiptRepo := repositories.NewReceiptRepository(client)
	receiptService := services.NewReceiptService(client, receiptRepo)
	receiptController := controllers.NewReceiptController(receiptService)

	router.POST("/receipts/process", receiptController.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptController.GetReceiptPoints)
}
