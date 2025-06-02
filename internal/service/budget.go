package service

import (
	"errors"
	"fmt"
	"time"
)

// Budget represents a monthly budget
type Budget struct {
	Month  time.Month `json:"month"`
	Year   int        `json:"year"`
	Amount float64    `json:"amount"`
}

// BudgetService handles budget-related operations
type BudgetService struct {
	expenseService *ExpenseService
	budgets        map[string]Budget // key is "year-month"
}

// NewBudgetService creates a new budget service
func NewBudgetService(expenseService *ExpenseService) *BudgetService {
	return &BudgetService{
		expenseService: expenseService,
		budgets:        make(map[string]Budget),
	}
}

// SetBudget sets a budget for a specific month and year
func (s *BudgetService) SetBudget(month int, year int, amount float64) error {
	if month < 1 || month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if amount <= 0 {
		return errors.New("budget amount must be greater than zero")
	}

	// If year is not specified (0), use current year
	if year == 0 {
		year = time.Now().Year()
	}

	key := fmt.Sprintf("%d-%d", year, month)
	s.budgets[key] = Budget{
		Month:  time.Month(month),
		Year:   year,
		Amount: amount,
	}

	return nil
}

// GetBudget returns the budget for a specific month and year
func (s *BudgetService) GetBudget(month int, year int) (Budget, error) {
	if month < 1 || month > 12 {
		return Budget{}, errors.New("month must be between 1 and 12")
	}

	// If year is not specified (0), use current year
	if year == 0 {
		year = time.Now().Year()
	}

	key := fmt.Sprintf("%d-%d", year, month)
	budget, exists := s.budgets[key]
	if !exists {
		return Budget{}, errors.New("budget not found for specified month and year")
	}

	return budget, nil
}

// CheckBudget checks if the current expenses exceed the budget for a specific month and year
// Returns remaining budget amount and a boolean indicating if budget is exceeded
func (s *BudgetService) CheckBudget(month int, year int) (float64, bool, error) {
	budget, err := s.GetBudget(month, year)
	if err != nil {
		return 0, false, err
	}

	summary, err := s.expenseService.GetMonthlySummary(month)
	if err != nil {
		return 0, false, err
	}

	remaining := budget.Amount - summary.TotalAmount
	exceeded := remaining < 0

	return remaining, exceeded, nil
}
