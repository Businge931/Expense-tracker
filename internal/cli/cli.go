package cli

import (
	"flag"
	"fmt"
	"time"

	"github.com/Businge931/expense-tracker/internal/models"
	"github.com/Businge931/expense-tracker/internal/service"
)

// CLI represents the command-line interface for the expense tracker
type CLI struct {
	expenseService *service.ExpenseService
	budgetService  *service.BudgetService
	exportService  *service.ExportService
}

// NewCLI creates a new CLI instance
func NewCLI(expenseService *service.ExpenseService, budgetService *service.BudgetService, exportService *service.ExportService) *CLI {
	return &CLI{
		expenseService: expenseService,
		budgetService:  budgetService,
		exportService:  exportService,
	}
}

// Run executes the CLI with the given arguments
func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		c.printUsage()
		return nil
	}

	command := args[0]

	switch command {
	case "add":
		return c.handleAddCommand(args[1:])
	case "list":
		return c.handleListCommand(args[1:])
	case "delete":
		return c.handleDeleteCommand(args[1:])
	case "summary":
		return c.handleSummaryCommand(args[1:])
	case "budget":
		return c.handleBudgetCommand(args[1:])
	case "export":
		return c.handleExportCommand(args[1:])
	case "help":
		c.printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// printUsage displays the usage information
func (c *CLI) printUsage() {
	fmt.Println("Expense Tracker - A simple tool to track your expenses")
	fmt.Println("\nUsage:")
	fmt.Println("  expense-tracker [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  add         Add a new expense")
	fmt.Println("  list        List all expenses")
	fmt.Println("  delete      Delete an expense")
	fmt.Println("  summary     Show a summary of expenses")
	fmt.Println("  budget      Set or check budget for a month")
	fmt.Println("  export      Export expenses to a CSV file")
	fmt.Println("  help        Show this help message")
	fmt.Println("\nOptions:")
	fmt.Println("  Run 'expense-tracker [command] --help' for command-specific help")
}

// handleAddCommand handles the 'add' command
func (c *CLI) handleAddCommand(args []string) error {
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("Usage: expense-tracker add --description DESCRIPTION --amount AMOUNT [--category CATEGORY]")
		return nil
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	description := addCmd.String("description", "", "Description of the expense")
	amount := addCmd.Float64("amount", 0, "Amount spent")
	category := addCmd.String("category", "", "Category of the expense (optional)")

	if err := addCmd.Parse(args); err != nil {
		return err
	}

	if *description == "" {
		return fmt.Errorf("description is required")
	}
	if *amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	id, err := c.expenseService.AddExpense(*description, *amount, *category)
	if err != nil {
		return err
	}

	fmt.Printf("Expense added successfully (ID: %d)\n", id)
	return nil
}

// handleListCommand handles the 'list' command
func (c *CLI) handleListCommand(args []string) error {
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("Usage: expense-tracker list")
		return nil
	}

	expenses, err := c.expenseService.GetAllExpenses()
	if err != nil {
		return err
	}

	if len(expenses) == 0 {
		fmt.Println("No expenses found")
		return nil
	}

	fmt.Println("ID\tDate\t\tDescription\tAmount")
	for _, expense := range expenses {
		fmt.Printf("%d\t%s\t%s\t$%.2f\n",
			expense.ID,
			expense.Date.Format("2006-01-02"),
			expense.Description,
			expense.Amount)
	}

	return nil
}

// handleDeleteCommand handles the 'delete' command
func (c *CLI) handleDeleteCommand(args []string) error {
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("Usage: expense-tracker delete --id ID")
		return nil
	}

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	id := deleteCmd.Int("id", 0, "ID of the expense to delete")

	if err := deleteCmd.Parse(args); err != nil {
		return err
	}

	if *id <= 0 {
		return fmt.Errorf("valid expense ID is required")
	}

	if err := c.expenseService.DeleteExpense(*id); err != nil {
		return err
	}

	fmt.Println("Expense deleted successfully")
	return nil
}

// handleSummaryCommand handles the 'summary' command
func (c *CLI) handleSummaryCommand(args []string) error {
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("Usage: expense-tracker summary [--month MONTH]")
		return nil
	}

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	month := summaryCmd.Int("month", 0, "Month to show summary for (1-12)")

	if err := summaryCmd.Parse(args); err != nil {
		return err
	}

	var summary interface{}
	var err error

	if *month > 0 {
		// Check if there's a budget for this month
		budget, budgetErr := c.budgetService.GetBudget(*month, 0)
		monthlySummary, err := c.expenseService.GetMonthlySummary(*month)
		if err != nil {
			return err
		}

		fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(*month).String(), monthlySummary.TotalAmount)

		// Show budget information if available
		if budgetErr == nil {
			remaining := budget.Amount - monthlySummary.TotalAmount
			fmt.Printf("Budget: $%.2f\n", budget.Amount)
			fmt.Printf("Remaining: $%.2f\n", remaining)

			if remaining < 0 {
				fmt.Printf("Warning: You've exceeded your budget by $%.2f\n", -remaining)
			}
		}

		return nil
	} else {
		summary, err = c.expenseService.GetExpenseSummary()
		if err != nil {
			return err
		}

		fmt.Printf("Total expenses: $%.2f\n", summary.(models.ExpenseSummary).TotalAmount)
		return nil
	}
}

// handleBudgetCommand handles the 'budget' command
func (c *CLI) handleBudgetCommand(args []string) error {
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("Usage: expense-tracker budget --month MONTH --amount AMOUNT")
		return nil
	}

	budgetCmd := flag.NewFlagSet("budget", flag.ExitOnError)
	month := budgetCmd.Int("month", 0, "Month to set budget for (1-12)")
	amount := budgetCmd.Float64("amount", 0, "Budget amount")

	if err := budgetCmd.Parse(args); err != nil {
		return err
	}

	if *month < 1 || *month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}
	if *amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	if err := c.budgetService.SetBudget(*month, 0, *amount); err != nil {
		return err
	}

	fmt.Printf("Budget of $%.2f set for %s\n", *amount, time.Month(*month).String())
	return nil
}

// handleExportCommand handles the 'export' command
func (c *CLI) handleExportCommand(args []string) error {
	if len(args) > 0 && args[0] == "--help" {
		fmt.Println("Usage: expense-tracker export --file FILEPATH")
		return nil
	}

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	file := exportCmd.String("file", "expenses.csv", "Path to export file")

	if err := exportCmd.Parse(args); err != nil {
		return err
	}

	if *file == "" {
		return fmt.Errorf("file path is required")
	}

	if err := c.exportService.ExportToCSV(*file); err != nil {
		return err
	}

	fmt.Printf("Expenses exported to %s\n", *file)
	return nil
}
