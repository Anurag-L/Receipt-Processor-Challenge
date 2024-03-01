package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var receiptData map[string]Receipt

func init() {
	receiptData = make(map[string]Receipt)
}

func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := GenerateReceiptID()
	receiptData[id] = receipt
	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]
	points := CalculatePoints(receiptID)

	if points == 0 {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := points
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GenerateReceiptID() string {
	id := uuid.New()
	return id.String()
}
func CalculatePoints(receiptID string) int {
	receipt := receiptData[receiptID]
	fmt.Printf("Receipt Data for ID %+v\n", receipt)
	points := 0
	// Calculate points for the retailer name (Rule 1)
	trimmedRetailer := strings.ReplaceAll(receipt.Retailer, " ", "")
	pointsRule1 := 0

	for _, char := range trimmedRetailer {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			pointsRule1++
		}
	}

	points += pointsRule1
	fmt.Printf("Rule 1:  %d points, total: %d points\n", pointsRule1, points)

	total, _ := strconv.ParseFloat(receipt.Total, 64)
	fmt.Printf("Receipt total in int %f", total)

	if math.Mod(total, 1.0) == 0 {
		points += 50
	}
	fmt.Printf("Rule 2: 50 points if added total: %d points\n", points)

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if total > 0.0 && math.Mod(total, 0.25) == 0.0 {
		points += 25
	}
	fmt.Printf("Rule 3: 25 points if added total: %d points\n", points)

	// Rule 4: 5 points for every two items on the receipt.
	pointsRule4 := 5 * (len(receipt.Items) / 2)
	points += pointsRule4
	fmt.Printf("Rule 4: 5 points for every two items on the receipt added: %d points total: %d points\n\n", pointsRule4, points)
	// Rule 5: Calculate points for item descriptions.
	for _, item := range receipt.Items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription))
		fmt.Printf("length of desc in this item %d, desc: %s\n", descriptionLength, item.ShortDescription)
		// If the trimmed length of the item description is a multiple of 3, calculate points.
		if descriptionLength%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			descriptionPoints := int(math.Ceil(price * 0.2))
			fmt.Printf("Amount of items %d\n, added in this run %d points\n", descriptionLength, descriptionPoints)
			points += descriptionPoints
		}

	}
	fmt.Printf("Rule 5: len item desc multiple of 3, total: %d points\n", points)

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)

	if purchaseDate.Day()%2 != 0 {
		points += 6
	}
	fmt.Printf("Rule 6: 6 points if added if purchase date odd, total: %d points\n", points)
	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.After(time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)) {
		points += 10
	}
	fmt.Printf("Rule 7: 10 points if added, total: %d points\n", points)
	fmt.Printf("Total points %d\n", points)
	return points
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", ProcessReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPointsHandler).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
