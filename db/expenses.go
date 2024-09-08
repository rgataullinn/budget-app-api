package db

import (
	"context"
	"personal-finance-api/models"
)

func AddExpenseInDb(newExpense models.Expense) error {
	sqlScript := `
		INSERT INTO expenses (user_id, category_id, amount,name, description, expense_date, expense_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := Pool.Exec(context.Background(), sqlScript,
		newExpense.User_id,
		newExpense.Category_id,
		newExpense.Amount,
		newExpense.Name,
		newExpense.Description,
		newExpense.Expense_date,
		newExpense.Expense_time,
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateExpenseInDb(expense models.Expense) error {
	sqlScript := `
		UPDATE expenses
		SET
			user_id = $2,
			category_id = $3,
			amount = $4,
			name = $5,
			description = $6,
			expense_date = $7,
			expense_time = $8
		WHERE 
			id = $1;
	`
	_, err := Pool.Exec(context.Background(), sqlScript,
		expense.Id,
		expense.User_id,
		expense.Category_id,
		expense.Amount,
		expense.Name,
		expense.Description,
		expense.Expense_date,
		expense.Expense_time,
	)

	if err != nil {
		return err
	}
	return nil
}

func GetExpenseFromDb(id int) (*models.Expense, error) {
	sqlScript := `
        SELECT id, user_id, category_id, amount,name, description, expense_date, expense_time
        FROM expenses
        WHERE id = $1
    `

	// Prepare a variable to hold the result
	var expense models.Expense

	// Execute the query
	row := Pool.QueryRow(context.Background(), sqlScript, id)

	// Scan the result into the expense struct
	err := row.Scan(&expense.Id, &expense.User_id, &expense.Category_id,
		&expense.Amount, &expense.Name, &expense.Description, &expense.Expense_date, &expense.Expense_time)
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

func GetAllExpensesByDate() ([]struct {
	Day      string           `json:"day"`
	Expenses []models.Expense `json:"expenses"`
}, error) {
	sqlScript := `
        SELECT id, user_id, category_id, amount,name, description, expense_date, expense_time
        FROM expenses
        ORDER BY expense_date
    `

	expensesByDate := make(map[string][]models.Expense)

	rows, err := Pool.Query(context.Background(), sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense

		err := rows.Scan(&expense.Id, &expense.User_id, &expense.Category_id,
			&expense.Amount, &expense.Name, &expense.Description, &expense.Expense_date, &expense.Expense_time)
		if err != nil {
			return nil, err
		}

		// Group expenses by date
		expensesByDate[expense.Expense_date] = append(expensesByDate[expense.Expense_date], expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	var result []struct {
		Day      string           `json:"day"`
		Expenses []models.Expense `json:"expenses"`
	}
	for day, expenses := range expensesByDate {
		result = append(result, struct {
			Day      string           `json:"day"`
			Expenses []models.Expense `json:"expenses"`
		}{
			Day:      day,
			Expenses: expenses,
		})
	}

	return result, nil
}

func GetAllExpenseByCategory() ([]struct {
	Category string           `json:"category"`
	Expenses []models.Expense `json:"expenses"`
}, error) {
	sqlScript := `
		SELECT id, user_id, category_id, amount,name, description, expense_date, expense_time
        FROM expenses
        ORDER BY category_id
	`
	expensesByCategory := make(map[string][]models.Expense)

	rows, err := Pool.Query(context.Background(), sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense

		err := rows.Scan(&expense.Id, &expense.User_id, &expense.Category_id,
			&expense.Amount, &expense.Name, &expense.Description, &expense.Expense_date, &expense.Expense_time)
		if err != nil {
			return nil, err
		}

		categoryName, err := GetCategoryName(expense.Category_id)
		if err != nil {
			return nil, err
		}

		expensesByCategory[categoryName] =
			append(expensesByCategory[categoryName], expense)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	var result []struct {
		Category string           `json:"category"`
		Expenses []models.Expense `json:"expenses"`
	}
	for category, expenses := range expensesByCategory {
		result = append(result, struct {
			Category string           `json:"category"`
			Expenses []models.Expense `json:"expenses"`
		}{
			Category: category,
			Expenses: expenses,
		})
	}

	return result, nil
}

func DeleteExpenseFromDb(id int) error {
	sqlScript := `
		DELETE FROM expenses
		WHERE id = $1;
	`
	_, err := Pool.Exec(context.Background(), sqlScript, id)
	if err != nil {
		return err
	}
	return nil
}
