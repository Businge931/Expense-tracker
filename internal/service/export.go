// internal/service/export.go
package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ExportService handles exporting expense data
type ExportService struct {
	expenseService *ExpenseService
}

// NewExportService creates a new export service
func NewExportService(expenseService *ExpenseService) *ExportService {
	return &ExportService{
		expenseService: expenseService,
	}
}

// ExportToCSV exports all expenses to a CSV file
func (s *ExportService) ExportToCSV(filePath string) error {
	expenses, err := s.expenseService.GetAllExpenses()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Date", "Description", "Amount", "Category"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	// Write data
	for _, expense := range expenses {
		record := []string{
			strconv.Itoa(expense.ID),
			expense.Date.Format("2006-01-02"),
			expense.Description,
			fmt.Sprintf("%.2f", expense.Amount),
			expense.Category,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}
