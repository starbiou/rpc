package controllers

import (
	"billing/internal/core/services"
	"billing/internal/db/ent"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReceiptController struct {
	service *services.ReceiptService
}

func NewReceiptController(service *services.ReceiptService) *ReceiptController {
	return &ReceiptController{service: service}
}

func (c *ReceiptController) ProcessReceipt(ctx *gin.Context) {
	var receipt ent.Receipt
	if err := ctx.ShouldBindJSON(&receipt); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	id, err := c.service.AddReceipt(ctx, receipt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process receipt: " +
			err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func (c *ReceiptController) GetReceiptPoints(ctx *gin.Context) {
	id := ctx.Param("id")

	points, err := c.service.GetReceiptPoints(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"points": points})
}
