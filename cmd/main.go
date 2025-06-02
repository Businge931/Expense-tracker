package main

import (
	"time"
)

type Expense struct {
	ID			int `json:"id"`
	Name		string `json:"name"`
	Amount		float64 `json:"amount"`
	Description	string `json:"description"`
	Category	string `json:"category"`
	CreatedAt	time.Time `json:"date"`
	UpdatedAt	time.Time `json:"-"`
}

var expensesFilePath string = "expenses.json"

func main() {
	// 1. Define the args

	// 2. Check if file exists

	// 3. Add new expense

	// 4. Save expenses to file

	// 5. Print expenses

}

