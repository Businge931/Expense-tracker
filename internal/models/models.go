package models

import (
	"fmt"
	"time"
)

type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category,omitempty"` // Optional for basic functionality
	Date        time.Time `json:"date"`
}

// String returns a formatted string representation of the expense
func (e Expense) String() string {
	return fmt.Sprintf("%d\t%s\t%s\t$%.2f",
		e.ID,
		e.Date.Format("2006-01-02"),
		e.Description,
		e.Amount)
}

type ExpenseSummary struct {
	TotalAmount    float64            `json:"totalAmount"`
	CategoryTotals map[string]float64 `json:"categoryTotals,omitempty"`
	ExpenseCount   int                `json:"expenseCount"`
	Month          time.Month         `json:"month,omitempty"`
	Year           int                `json:"year,omitempty"`
}

func (s ExpenseSummary) String() string {
	if s.Month > 0 {
		return fmt.Sprintf("Total expenses for %s: $%.2f", s.Month.String(), s.TotalAmount)
	}
	return fmt.Sprintf("Total expenses: $%.2f", s.TotalAmount)
}
