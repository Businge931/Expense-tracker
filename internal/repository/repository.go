package repository

import (
	"time"

	"github.com/Businge931/expense-tracker/internal/models"
)

type ExpenseRepository interface {
	Add(description string, amount float64, category string) (int, error)
	GetByID(id int) (models.Expense, error)
	GetAll() ([]models.Expense, error)
	GetByMonth(month time.Month, year int) ([]models.Expense, error)
	Delete(id int) error
	GetSummary() (models.ExpenseSummary, error)
	GetMonthlySummary(month time.Month, year int) (models.ExpenseSummary, error)
}
