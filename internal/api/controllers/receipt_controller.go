package controllers

import (
	"billing/internal/core/services"
	"billing/internal/db/ent"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ItemPayload struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required,containsany=\\d+\\.\\d{2}$"`
}

type ReceiptPayload struct {
	Retailer     string        `json:"retailer" validate:"required"`
	PurchaseDate string        `json:"purchaseDate" validate:"required"`
	PurchaseTime string        `json:"purchaseTime" validate:"required"`
	Total        string        `json:"total" validate:"required,containsany=\\d+\\.\\d{2}$"`
	Items        []ItemPayload `json:"items" validate:"required,dive"`
}

type ReceiptController struct {
	service services.ReceiptServicer
}

// parsePrice converts a price string to an integer representing cents
func parsePrice(priceStr string) (int, error) {
	priceStr = strings.TrimSpace(priceStr)
	return strconv.Atoi(strings.Replace(priceStr, ".", "", 1))
}

// convertItems converts a slice of ItemPayload to a slice of *ent.Item
func convertItems(itemPayloads []ItemPayload) ([]*ent.Item, error) {
	var items []*ent.Item
	for _, item := range itemPayloads {
		priceInt, err := parsePrice(item.Price)
		if err != nil {
			return nil, fmt.Errorf("invalid item price format: %w", err)
		}
		log.Printf("Adding item: %s, price: %d\n", item.ShortDescription, priceInt)
		items = append(items, &ent.Item{
			ShortDescription: strings.TrimSpace(item.ShortDescription),
			Price:            priceInt,
		})
	}
	return items, nil
}

func NewReceiptController(service services.ReceiptServicer) *ReceiptController {
	return &ReceiptController{service: service}
}

func (c *ReceiptController) ProcessReceipt(ctx *gin.Context) {
	payload, exists := ctx.Get("payload")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Payload not found"})
		return
	}

	receiptPayload, ok := payload.(ReceiptPayload)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload format"})
		return
	}

	totalInt, err := parsePrice(receiptPayload.Total)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please verify input. Invalid total format"})
		return
	}

	fmt.Printf("Processing receipt for retailer: %s, total: %d\n", receiptPayload.Retailer, totalInt)

	receipt := ent.Receipt{
		Retailer:     strings.TrimSpace(receiptPayload.Retailer),
		PurchaseDate: strings.TrimSpace(receiptPayload.PurchaseDate),
		PurchaseTime: strings.TrimSpace(receiptPayload.PurchaseTime),
		Total:        totalInt,
	}

	// Convert items
	items, err := convertItems(receiptPayload.Items)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	receipt.Edges.Items = items

	id, err := c.service.AddReceipt(ctx, receipt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process receipt: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (c *ReceiptController) GetReceiptPoints(ctx *gin.Context) {
	id := ctx.Param("id")
	receiptID, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt ID format"})
		return
	}

	points, err := c.service.GetReceiptPoints(ctx, receiptID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve receipt points: " + err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"points": points})
}
