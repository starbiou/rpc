package services

import (
	"billing/internal/db/ent"
	"billing/internal/db/repositories"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// RetailerNamePointsMultiplier is the multiplier for points based on retailer name characters.
	RetailerNamePointsMultiplier = 1
	// RoundDollarPoints is the points awarded for a round dollar total.
	RoundDollarPoints = 50
	// MultipleOfQuarterPoints is the points awarded if the total is a multiple of 0.25.
	MultipleOfQuarterPoints = 25
	// TwoItemsPoints is the points awarded for every two items on the receipt.
	TwoItemsPoints = 5
	// DescriptionLengthMultiplier is the multiplier for points based on item description length.
	DescriptionLengthMultiplier = 0.2
	// OddDayPoints is the points awarded if the purchase day is odd.
	OddDayPoints = 6
	// AfternoonPurchasePoints is the points awarded if the purchase time is between 2:00pm and 4:00pm.
	AfternoonPurchasePoints = 10
)

// ReceiptService provides methods to process receipts and calculate points.
type ReceiptService struct {
	client *ent.Client
	repo   repositories.ReceiptRepository
}

// NewReceiptService creates a new instance of ReceiptService with the given repository.
func NewReceiptService(client *ent.Client, repo repositories.ReceiptRepository) *ReceiptService {
	return &ReceiptService{client: client, repo: repo}
}

// AddReceipt adds a new receipt to the repository and returns its ID.
func (s *ReceiptService) AddReceipt(ctx context.Context, receipt ent.Receipt) (string, error) {
	createdReceipt, err := s.repo.CreateReceipt(ctx, &receipt)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(createdReceipt.ID), nil
}

// GetReceiptPoints retrieves a receipt by ID and calculates the points awarded.
func (s *ReceiptService) GetReceiptPoints(ctx context.Context, id string) (int, error) {
	receiptID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	receipt, err := s.repo.GetReceiptByID(ctx, receiptID)
	if err != nil {
		return 0, err
	}
	return s.CalculatePoints(receipt), nil
}

func parseTotal(totalStr string) (float64, error) {
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid total value: %v", err)
	}
	return total, nil
}

// CalculatePoints calculates the total points for a given receipt based on various rules.
func (s *ReceiptService) CalculatePoints(receipt *ent.Receipt) int {
	points := 0
	points += s.calculateRetailerNamePoints(receipt.Retailer)

	total, err := parseTotal(receipt.Total)
	if err == nil {
		points += s.calculateTotalPoints(total)
	} else {
		// Handle the error, e.g., log it or return a default value
	}

	points += s.calculateItemPoints(receipt.Edges.Items)
	points += s.calculateDatePoints(receipt.PurchaseDate)
	points += s.calculateTimePoints(receipt.PurchaseTime)
	return points
}

// Rule 1: Retailer name points
// calculateRetailerNamePoints calculates points based on the number of alphanumeric characters in the retailer's name.
func (s *ReceiptService) calculateRetailerNamePoints(retailer string) int {
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	return len(re.FindAllString(retailer, -1)) * RetailerNamePointsMultiplier
}

// calculateTotalPoints calculates points based on the total amount of the receipt.
func (s *ReceiptService) calculateTotalPoints(total float64) int {
	points := 0
	if total == float64(int(total)) { // Rule 2: Round dollar amount
		points += RoundDollarPoints
	}
	if int(total*100)%25 == 0 { // Rule 3: Multiple of 0.25
		points += MultipleOfQuarterPoints
	}
	return points
}

// calculateItemPoints calculates points based on the number of items and their descriptions.
func (s *ReceiptService) calculateItemPoints(items []*ent.Item) int {
	points := (len(items) / 2) * TwoItemsPoints
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 { // Rule 4: Every two items
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(priceFloat*DescriptionLengthMultiplier + 0.5) // Rule 5: Description length multiple of 3
			}
		}
	}
	return points
}

// Rule 6: Odd day
// calculateDatePoints calculates points based on the purchase date.
func (s *ReceiptService) calculateDatePoints(date string) int {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}
	if parsedDate.Day()%2 != 0 {
		return OddDayPoints
	}
	return 0
}

// Rule 7: Time between 2:00pm and 4:00pm
// calculateTimePoints calculates points based on the purchase time.
func (s *ReceiptService) calculateTimePoints(timeStr string) int {
	hour, err := strconv.Atoi(strings.Split(timeStr, ":")[0])
	if err != nil {
		return 0
	}
	if hour >= 14 && hour < 16 {
		return AfternoonPurchasePoints
	}
	return 0
}
