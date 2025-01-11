package routes

import (
	"receipt-processor/api/controllers"
	"receipt-processor/api/services"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	service := services.NewReceiptService()
	controller := controllers.NewReceiptController(service)

	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", controller.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", controller.GetReceiptPoints).Methods("GET")

	return router
}
