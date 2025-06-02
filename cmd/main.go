package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Businge931/expense-tracker/internal/cli"
	"github.com/Businge931/expense-tracker/internal/repository"
	"github.com/Businge931/expense-tracker/internal/service"
)

func main() {
	// Set up the data directory in the project
	execDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	dataDir := filepath.Join(execDir, "data")

	// Initialize repository
	repo, err := repository.NewJSONFileRepository(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing repository: %v\n", err)
		os.Exit(1)
	}

	// Initialize services
	expenseService := service.NewExpenseService(repo)
	budgetService := service.NewBudgetService(expenseService)
	exportService := service.NewExportService(expenseService)

	// Initialize CLI
	cli := cli.NewCLI(expenseService, budgetService, exportService)

	// Run CLI with command-line arguments
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
