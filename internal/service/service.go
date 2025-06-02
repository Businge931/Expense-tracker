package service

import (
	"errors"
	"time"

	"github.com/Businge931/expense-tracker/internal/models"
	"github.com/Businge931/expense-tracker/internal/repository"
)

// ExpenseService handles business logic for expense operations
type ExpenseService struct {
	repo repository.ExpenseRepository
}

// NewExpenseService creates a new expense service
func NewExpenseService(repo repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{
		repo: repo,
	}
}

// AddExpense adds a new expense
func (s *ExpenseService) AddExpense(description string, amount float64, category string) (int, error) {
	// Validate inputs
	if description == "" {
		return 0, errors.New("description cannot be empty")
	}
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}

	// Add expense to repository
	return s.repo.Add(description, amount, category)
}

// GetAllExpenses returns all expenses
func (s *ExpenseService) GetAllExpenses() ([]models.Expense, error) {
	return s.repo.GetAll()
}

// GetExpenseByID returns an expense with the given ID
func (s *ExpenseService) GetExpenseByID(id int) (models.Expense, error) {
	if id <= 0 {
		return models.Expense{}, errors.New("invalid expense ID")
	}
	return s.repo.GetByID(id)
}

// DeleteExpense deletes an expense with the given ID
func (s *ExpenseService) DeleteExpense(id int) error {
	if id <= 0 {
		return errors.New("invalid expense ID")
	}
	return s.repo.Delete(id)
}

// GetExpenseSummary returns a summary of all expenses
func (s *ExpenseService) GetExpenseSummary() (models.ExpenseSummary, error) {
	return s.repo.GetSummary()
}

// GetMonthlySummary returns a summary of expenses for a specific month
func (s *ExpenseService) GetMonthlySummary(month int) (models.ExpenseSummary, error) {
	if month < 1 || month > 12 {
		return models.ExpenseSummary{}, errors.New("month must be between 1 and 12")
	}

	// Use the current year if not specified
	currentYear := time.Now().Year()
	return s.repo.GetMonthlySummary(time.Month(month), currentYear)
}
