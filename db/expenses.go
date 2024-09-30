package db

import (
	"context"
	"personal-finance-api/models"
	"sort"
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

func GetAllExpensesByDate(month int) ([]struct {
	Day      string           `json:"day"`
	Expenses []models.Expense `json:"expenses"`
}, error) {
	// TODO: organize the same as by category
	sqlScript := `
    SELECT e.id, e.user_id, e.category_id, c.name AS category, c.color, e.amount, e.name, e.description, e.expense_date, e.expense_time
    FROM expenses e
    JOIN categories c ON e.category_id = c.id
	WHERE EXTRACT(MONTH FROM expense_date::date) = $1
    ORDER BY e.expense_date DESC
`

	expensesByDate := make(map[string][]models.Expense)

	rows, err := Pool.Query(context.Background(), sqlScript, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense

		err := rows.Scan(&expense.Id, &expense.User_id, &expense.Category_id, &expense.Category, &expense.Color,
			&expense.Amount, &expense.Name, &expense.Description, &expense.Expense_date, &expense.Expense_time)
		if err != nil {
			return nil, err
		}

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

		sort.Slice(expenses, func(i, j int) bool {
			return expenses[i].Expense_time > expenses[j].Expense_time
		})

		result = append(result, struct {
			Day      string           `json:"day"`
			Expenses []models.Expense `json:"expenses"`
		}{
			Day:      day,
			Expenses: expenses,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Day > result[j].Day
	})

	return result, nil
}

func GetAllExpenseByCategory(month int) ([]struct {
	Category string           `json:"category"`
	Total    float32          `json:"total"`
	Color    string           `json:"color"`
	Expenses []models.Expense `json:"expenses"`
}, error) {
	sqlScript := `
			SELECT c.id, c.name, c.color, SUM(e.amount) AS total
			FROM categories c
			LEFT JOIN expenses e
			ON c.id = e.category_id 
			WHERE EXTRACT(MONTH FROM e.expense_date::date) = $1
			GROUP BY c.id
			ORDER BY total DESC;
		`
	rows, err := Pool.Query(context.Background(), sqlScript, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []struct {
		Category string           `json:"category"`
		Total    float32          `json:"total"`
		Color    string           `json:"color"`
		Expenses []models.Expense `json:"expenses"`
	}

	for rows.Next() {
		var sub_res struct {
			Category string           `json:"category"`
			Total    float32          `json:"total"`
			Color    string           `json:"color"`
			Expenses []models.Expense `json:"expenses"`
		}
		var category_id int
		err := rows.Scan(&category_id, &sub_res.Category, &sub_res.Color, &sub_res.Total)
		if err != nil {
			return nil, err
		}

		sub_res.Expenses, err = GetExpensesByCategory(category_id)
		if err != nil {
			return nil, err
		}
		result = append(result, sub_res)
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
