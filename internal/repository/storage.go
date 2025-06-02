package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/Businge931/expense-tracker/internal/models"
)

// JSONFileRepository implements ExpenseRepository using a JSON file for storage
type JSONFileRepository struct {
	filePath string
	mutex    sync.RWMutex
}

// NewJSONFileRepository creates a new repository that stores data in a JSON file
func NewJSONFileRepository(dataDir string) (*JSONFileRepository, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	filePath := filepath.Join(dataDir, "expenses.json")

	// Create file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Initialize with empty expenses array
		initialData := struct {
			Expenses []models.Expense `json:"expenses"`
		}{
			Expenses: []models.Expense{},
		}

		data, err := json.MarshalIndent(initialData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal initial data: %w", err)
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to create initial expenses file: %w", err)
		}
	}

	return &JSONFileRepository{
		filePath: filePath,
	}, nil
}

// loadExpenses reads all expenses from the JSON file
func (r *JSONFileRepository) loadExpenses() ([]models.Expense, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read expenses file: %w", err)
	}

	var data struct {
		Expenses []models.Expense `json:"expenses"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expenses: %w", err)
	}

	return data.Expenses, nil
}

// saveExpenses writes all expenses to the JSON file
func (r *JSONFileRepository) saveExpenses(expenses []models.Expense) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	data := struct {
		Expenses []models.Expense `json:"expenses"`
	}{
		Expenses: expenses,
	}

	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal expenses: %w", err)
	}

	if err := os.WriteFile(r.filePath, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write expenses file: %w", err)
	}

	return nil
}

// Add adds a new expense and returns its ID
func (r *JSONFileRepository) Add(description string, amount float64, category string) (int, error) {
	expenses, err := r.loadExpenses()
	if err != nil {
		return 0, err
	}

	// Find the highest ID and increment
	maxID := 0
	for _, e := range expenses {
		if e.ID > maxID {
			maxID = e.ID
		}
	}

	// Create new expense
	expense := models.Expense{
		ID:          maxID + 1,
		Description: description,
		Amount:      amount,
		Category:    category,
		Date:        time.Now(),
	}

	expenses = append(expenses, expense)

	if err := r.saveExpenses(expenses); err != nil {
		return 0, err
	}

	return expense.ID, nil
}

// GetByID retrieves an expense by its ID
func (r *JSONFileRepository) GetByID(id int) (models.Expense, error) {
	expenses, err := r.loadExpenses()
	if err != nil {
		return models.Expense{}, err
	}

	for _, expense := range expenses {
		if expense.ID == id {
			return expense, nil
		}
	}

	return models.Expense{}, errors.New("expense not found")
}

// GetAll retrieves all expenses
func (r *JSONFileRepository) GetAll() ([]models.Expense, error) {
	return r.loadExpenses()
}

// GetByMonth retrieves expenses for a specific month and year
func (r *JSONFileRepository) GetByMonth(month time.Month, year int) ([]models.Expense, error) {
	expenses, err := r.loadExpenses()
	if err != nil {
		return nil, err
	}

	var result []models.Expense
	for _, expense := range expenses {
		if expense.Date.Month() == month && expense.Date.Year() == year {
			result = append(result, expense)
		}
	}

	return result, nil
}

// Delete removes an expense by its ID
func (r *JSONFileRepository) Delete(id int) error {
	expenses, err := r.loadExpenses()
	if err != nil {
		return err
	}

	foundIndex := -1
	for i, expense := range expenses {
		if expense.ID == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return errors.New("expense not found")
	}

	// Remove the expense from the slice
	expenses = slices.Delete(expenses, foundIndex, foundIndex+1)

	return r.saveExpenses(expenses)
}

// GetSummary returns a summary of all expenses
func (r *JSONFileRepository) GetSummary() (models.ExpenseSummary, error) {
	expenses, err := r.loadExpenses()
	if err != nil {
		return models.ExpenseSummary{}, err
	}

	summary := models.ExpenseSummary{
		TotalAmount:    0,
		CategoryTotals: make(map[string]float64),
		ExpenseCount:   len(expenses),
	}

	for _, expense := range expenses {
		summary.TotalAmount += expense.Amount
		if expense.Category != "" {
			summary.CategoryTotals[expense.Category] += expense.Amount
		}
	}

	return summary, nil
}

// GetMonthlySummary returns a summary of expenses for a specific month and year
func (r *JSONFileRepository) GetMonthlySummary(month time.Month, year int) (models.ExpenseSummary, error) {
	expenses, err := r.GetByMonth(month, year)
	if err != nil {
		return models.ExpenseSummary{}, err
	}

	summary := models.ExpenseSummary{
		TotalAmount:    0,
		CategoryTotals: make(map[string]float64),
		ExpenseCount:   len(expenses),
		Month:          month,
		Year:           year,
	}

	for _, expense := range expenses {
		summary.TotalAmount += expense.Amount
		if expense.Category != "" {
			summary.CategoryTotals[expense.Category] += expense.Amount
		}
	}

	return summary, nil
}
