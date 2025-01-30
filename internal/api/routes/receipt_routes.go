package routes

import (
	"billing/internal/api/controllers"
	"billing/internal/api/middleware"
	"billing/internal/core/services"
	"billing/internal/db/repositories"
)

func SetupReceiptRoutes(config RouteConfig) {
	receiptRepo := repositories.NewReceiptRepository(config.Client)
	receiptService := services.NewReceiptService(config.Client, receiptRepo)
	receiptController := controllers.NewReceiptController(receiptService)

	config.Router.POST("/receipts/process",
		middleware.ValidationMiddleware[controllers.ReceiptPayload](config.Validator),
		receiptController.ProcessReceipt)

	config.Router.GET("/receipts/:id/points",
		receiptController.GetReceiptPoints)
}
