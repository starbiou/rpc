package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"receipt-processor/api/services"
	"receipt-processor/db/schemas"

	"github.com/gorilla/mux"
)

type ReceiptController struct {
	service *services.ReceiptService
}

func NewReceiptController(service *services.ReceiptService) *ReceiptController {
	return &ReceiptController{service: service}
}

func (c *ReceiptController) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt schemas.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := c.service.AddReceipt(context.Background(), receipt)
	if err != nil {
		http.Error(w, "Failed to process receipt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (c *ReceiptController) GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, err := c.service.GetReceiptPoints(context.Background(), id)
	if err != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}
