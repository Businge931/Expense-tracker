# Expense Tracker

A simple command-line expense tracker application to manage your finances. This application allows you to add, delete, and view your expenses, as well as generate summaries and export data.

Project based on: [roadmap.sh/projects/expense-tracker](https://roadmap.sh/projects/expense-tracker)

## Features

- Add expenses with description and amount
- Add optional category to expenses
- Update existing expenses
- Delete expenses
- View all expenses in a tabular format
- View summary of all expenses
- View monthly expense summaries
- Set and track monthly budgets
- Export expenses to CSV

## Installation

Clone the repository and build the application:

```bash
git clone https://github.com/Businge931/expense-tracker.git
cd expense-tracker
go build -o expense-tracker ./cmd/main.go
```

## Usage

### Help

Display all available commands:

```bash
./expense-tracker help
```

Get help for a specific command:

```bash
./expense-tracker [command] --help
```

### Adding Expenses

Add a new expense with required description and amount:

```bash
./expense-tracker add --description "Lunch" --amount 20
```

Add with optional category:

```bash
./expense-tracker add --description "Groceries" --amount 50 --category "Food"
```

### Viewing Expenses

List all expenses:

```bash
./expense-tracker list
```

### Deleting Expenses

Delete an expense by ID:

```bash
./expense-tracker delete --id 1
```

### Expense Summaries

View summary of all expenses:

```bash
./expense-tracker summary
```

View summary for a specific month (1-12):

```bash
./expense-tracker summary --month 6
```

### Budget Management

Set a budget for a specific month:

```bash
./expense-tracker budget --month 6 --amount 1000
```

The summary command will show budget status when a budget exists for the specified month.

### Exporting Data

Export all expenses to a CSV file:

```bash
./expense-tracker export --file expenses.csv
```

## Data Storage

Expense data is stored in a JSON file located in the `data` directory:

```
data/expenses.json
```

## Examples

Here are some examples of the commands in action:

```bash
# Add some expenses
./expense-tracker add --description "Lunch" --amount 20
# Expense added successfully (ID: 1)

./expense-tracker add --description "Dinner" --amount 10
# Expense added successfully (ID: 2)

# List all expenses
./expense-tracker list
# ID  Date       Description  Amount
# 1   2025-06-02  Lunch        $20
# 2   2025-06-02  Dinner       $10

# View summary
./expense-tracker summary
# Total expenses: $30

# Delete an expense
./expense-tracker delete --id 2
# Expense deleted successfully

# View updated summary
./expense-tracker summary
# Total expenses: $20

# View monthly summary
./expense-tracker summary --month 6
# Total expenses for June: $20
```

## License

This project is open source and available under the [MIT License](LICENSE).
