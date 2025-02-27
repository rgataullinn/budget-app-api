package db

import (
	"context"
	"personal-finance-api/models"
)

func CreateExpense(newExpense models.Expense) error {
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

func UpdateExpense(expense models.Expense) error {
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

func GetExpense(id int) (models.Expense, error) {
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
		return models.Expense{}, err
	}

	return expense, nil
}

func GetAllExpensesGroupedByCategory(month int, user_id int) (
	[]models.ExpensesGroupedByCategory, error) {
	sqlScript := `
			SELECT c.id, c.name
			FROM categories c;
	`
	rows, err := Pool.Query(context.Background(), sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.Id, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	var result []models.ExpensesGroupedByCategory

	for _, c := range categories {
		expenses, total, err := GetAllExpensesByCategoryId(c.Id, user_id, month)
		if err != nil {
			return nil, err
		}
		if len(expenses) == 0 {
			continue
		}
		var sub_res models.ExpensesGroupedByCategory
		sub_res.Category = c
		sub_res.Total = total
		sub_res.Expenses = expenses
		result = append(result, sub_res)
	}
	return result, nil
}

func GetAllExpensesByCategoryId(category_id int, user_id int, month int) (
	[]models.Expense, float32, error) {
	sqlScript := `
			SELECT e.id, e.user_id, e.category_id, 
				e.amount, e.name, e.expense_date, e.expense_time
			FROM expenses e
			WHERE e.category_id = $1 and e.user_id = $2 and 
				EXTRACT(MONTH FROM e.expense_date::date) = $3
		`
	rows, err := Pool.Query(context.Background(), sqlScript, category_id, user_id, month)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []models.Expense

	var totalAmount float32

	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.Id, &e.User_id, &e.Category_id, &e.Amount,
			&e.Name, &e.Expense_date, &e.Expense_time)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, e)
		totalAmount += float32(e.Amount)
	}
	return expenses, totalAmount, nil
}

func GetAllExpensesGroupedByDay(month int, user_id int) (
	[]models.ExpensesGroupedByDay, error) {
	sqlScript := `
		select DISTINCT e.expense_date
		FROM expenses e
		WHERE EXTRACT(MONTH FROM e.expense_date::date) = $1 and
			e.user_id = $2
		ORDER BY e.expense_date asc;
	`
	rows, err := Pool.Query(context.Background(), sqlScript, month, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var days []string

	for rows.Next() {
		var day string
		err := rows.Scan(&day)
		if err != nil {
			return nil, err
		}
		days = append(days, day)
	}

	var result []models.ExpensesGroupedByDay

	for _, day := range days {
		expenses, total, err := GetAllExpensesByDay(day, user_id, month)
		if err != nil {
			return nil, err
		}
		var sub_res models.ExpensesGroupedByDay
		sub_res.Day = day
		sub_res.Total = total
		sub_res.Expenses = expenses

		result = append(result, sub_res)
	}
	return result, nil
}

func GetAllExpensesByDay(day string, user_id int, month int) (
	[]models.Expense, float32, error) {
	sqlScript := `
			SELECT e.id, e.user_id, e.category_id, c.name, 
				e.amount, e.name, e.expense_date, e.expense_time
			FROM expenses e
			JOIN categories c
			ON e.category_id = c.id
			WHERE e.expense_date = $1 and e.user_id = $2 and 
				EXTRACT(MONTH FROM e.expense_date::date) = $3
		`
	rows, err := Pool.Query(context.Background(), sqlScript, day, user_id, month)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []models.Expense

	var totalAmount float32

	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.Id, &e.User_id, &e.Category_id, &e.Category, &e.Amount,
			&e.Name, &e.Expense_date, &e.Expense_time)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, e)
		totalAmount += float32(e.Amount)
	}
	return expenses, totalAmount, nil
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
