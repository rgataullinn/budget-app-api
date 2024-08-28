package db

import (
	"context"
	"database/sql"
	"personal-finance-api/models"
)

func AddExpenseInDb(newExpense models.Expense) error {
	sqlScript := `
		INSERT INTO expenses (user_id, category_id, amount, description, expense_date)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := Pool.Exec(context.Background(), sqlScript,
		newExpense.User_id,
		newExpense.Category_id,
		newExpense.Amount,
		newExpense.Description,
		newExpense.Expense_date)
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
			description = $5,
			expense_date = $6
		WHERE 
			id = $1;
	`
	_, err := Pool.Exec(context.Background(), sqlScript,
		expense.Id,
		expense.User_id,
		expense.Category_id,
		expense.Amount,
		expense.Description,
		expense.Expense_date)

	if err != nil {
		return err
	}
	return nil
}

func GetExpenseFromDb(id int) (*models.Expense, error) {
	sqlScript := `
        SELECT id, user_id, category_id, amount, description, expense_date
        FROM expenses
        WHERE id = $1
    `

	// Prepare a variable to hold the result
	var expense models.Expense

	// Execute the query
	row := Pool.QueryRow(context.Background(), sqlScript, id)

	var date sql.NullString
	// Scan the result into the expense struct
	err := row.Scan(&expense.Id, &expense.User_id, &expense.Category_id,
		&expense.Amount, &expense.Description, &date)
	expense.Expense_date = date.String
	if err != nil {
		return nil, err
	}

	return &expense, nil
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
