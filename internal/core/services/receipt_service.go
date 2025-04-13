package services

import (
	"billing/internal/db/ent"
	"billing/internal/db/repositories"
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	RetailerNamePointsMultiplier = 1   // Multiplier for points based on retailer name characters
	RoundDollarPoints            = 50  // Points awarded for a round dollar total
	MultipleOfQuarterPoints      = 25  // Points awarded if the total is a multiple of 0.25
	TwoItemsPoints               = 5   // Points awarded for every two items on the receipt
	DescriptionLengthMultiplier  = 0.2 // Multiplier for points based on item description length
	OddDayPoints                 = 6   // Points awarded if the purchase day is odd
	AfternoonPurchasePoints      = 10  // Points awarded if the purchase time is between 2:00pm and 4:00pm
)

// ReceiptServicer defines the interface for receipt processing operations
type ReceiptServicer interface {
	AddReceipt(ctx context.Context, receipt ent.Receipt) (string, error)
	GetReceiptPoints(ctx context.Context, id int) (int, error)
}

// ReceiptService provides methods to process receipts and calculate points.
type ReceiptService struct {
	client      *ent.Client
	receiptRepo repositories.ReceiptRepository
}

// NewReceiptService creates a new instance of ReceiptService with the given repository.
func NewReceiptService(client *ent.Client, receiptRepo repositories.ReceiptRepository) *ReceiptService {
	return &ReceiptService{client: client, receiptRepo: receiptRepo}
}

// AddReceipt adds a new receipt, calculates points, and returns its ID.
func (s *ReceiptService) AddReceipt(ctx context.Context, receipt ent.Receipt) (string, error) {
	points := s.calculatePoints(&receipt)
	receipt.Points = points
	fmt.Println("receipt ", receipt)

	createdReceipt, err := s.receiptRepo.CreateReceiptWithItems(ctx, &receipt, receipt.Edges.Items)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(createdReceipt.ID), nil
}

// GetReceiptPoints retrieves the points for a receipt from the database.
func (s *ReceiptService) GetReceiptPoints(ctx context.Context, id int) (int, error) {
	receipt, err := s.receiptRepo.GetReceiptByID(ctx, id)
	if err != nil || receipt == nil {
		return 0, fmt.Errorf("receipt not found")
	}
	return receipt.Points, nil
}

// calculatePoints calculates the total points for a given receipt based on various rules.
func (s *ReceiptService) calculatePoints(receipt *ent.Receipt) int {
	points := 0
	points += s.calculateRetailerNamePoints(receipt.Retailer) // Rule 1: One point per alphanumeric character
	points += s.calculateTotalPoints(receipt.Total)           // Rules 2 & 3: Round dollar and quarter multiples
	points += s.calculateItemPoints(receipt.Edges.Items)      // Rules 4 & 5: Pairs of items and description length
	points += s.calculateDatePoints(receipt.PurchaseDate)     // Rule 6: Odd day bonus
	points += s.calculateTimePoints(receipt.PurchaseTime)     // Rule 7: Afternoon bonus
	return points
}

// calculateRetailerNamePoints calculates points based on the number of alphanumeric characters in the retailer's name.
func (s *ReceiptService) calculateRetailerNamePoints(retailer string) int {
	return len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(retailer, -1)) * RetailerNamePointsMultiplier
}

// calculateTotalPoints calculates points based on the total amount of the receipt.
func (s *ReceiptService) calculateTotalPoints(total int) int {
	points := 0
	if total%100 == 0 {
		points += RoundDollarPoints
	}
	if total%25 == 0 {
		points += MultipleOfQuarterPoints
	}
	return points
}

// calculateItemPoints calculates points based on the number of items and their descriptions.
func (s *ReceiptService) calculateItemPoints(items []*ent.Item) int {
	points := (len(items) / 2) * TwoItemsPoints
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			priceInDollars := float64(item.Price) / 100.0
			points += int(math.Ceil(priceInDollars * DescriptionLengthMultiplier))
		}
	}
	return points
}

// calculateDatePoints calculates points based on the purchase date.
func (s *ReceiptService) calculateDatePoints(date string) int {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Invalid date format: %v", err)
		return 0
	}
	if parsedDate.Day()%2 != 0 {
		return OddDayPoints
	}
	return 0
}

// calculateTimePoints calculates points based on the purchase time.
func (s *ReceiptService) calculateTimePoints(timeStr string) int {
	hour, err := strconv.Atoi(strings.Split(timeStr, ":")[0])
	if err != nil {
		log.Printf("Invalid time format: %v", err)
		return 0
	}
	if hour >= 14 && hour < 16 {
		return AfternoonPurchasePoints
	}
	return 0
}
