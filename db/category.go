package db

import (
	"context"
	"personal-finance-api/models"
)

func AddCategoryInDb(newCategory models.Category) error {
	sqlScript := `
		INSERT INTO categories (name, description)
		VALUES($1, $2);
	`
	_, err := Pool.Exec(context.Background(), sqlScript,
		newCategory.Name,
		newCategory.Description)

	return err
}

func UpdateCategoryInDb(newCategory models.Category) error {
	sqlScript := `
		UPDATE categories
		SET
			name = $2,
			description = $3
		WHERE
			id = $1
`
	_, err := Pool.Exec(context.Background(), sqlScript,
		newCategory.Id,
		newCategory.Name,
		newCategory.Description)

	return err
}

func GetCategoryInDb(id int) (models.Category, error) {
	sqlScript := `
		SELECT id, name, description, color
		FROM categories
		WHERE id= $1
	`
	var category models.Category
	row := Pool.QueryRow(context.Background(), sqlScript, id)
	err := row.Scan(&category.Id, &category.Name, &category.Description, &category.Color)
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func DeleteCategoryInDb(id int) error {
	sqlScript := `
		DELETE FROM categories
		WHERE id = $1
	`
	_, err := Pool.Exec(context.Background(), sqlScript, id)
	return err
}

func GetCategoryName(id int) (string, error) {
	sqlScript := `
		SELECT name FROM categories
		WHERE id=$1
	`
	var name string
	err := Pool.QueryRow(context.Background(), sqlScript, id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func GetAllCategoriesFromDb() ([]struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Id    int    `json:"id"`
}, error) {
	sqlScript := `
		SELECT id, name, color
		FROM categories
	`
	var result []struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		Id    int    `json:"id"`
	}
	rows, err := Pool.Query(context.Background(), sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var id int
		var color string
		err := rows.Scan(&id, &name, &color)
		if err != nil {
			return nil, err
		}
		result = append(result, struct {
			Name  string `json:"name"`
			Color string `json:"color"`
			Id    int    `json:"id"`
		}{
			Name:  name,
			Color: color,
			Id:    id,
		})
	}
	return result, nil
}

func GetAllCategoriesWithTotals() ([]struct {
	Name  string  `json:"name"`
	Color string  `json:"color"`
	Total float32 `json:"total"`
}, error) {
	sqlScript := `
		select c.name, sum(e.amount) as total, c.color
		from categories c 
		join expenses e 
		on c.id = e.category_id 
		group by c.id;
	`
	var result []struct {
		Name  string  `json:"name"`
		Color string  `json:"color"`
		Total float32 `json:"total"`
	}
	rows, err := Pool.Query(context.Background(), sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var total float32
		var color string
		err := rows.Scan(&name, &total, &color)
		if err != nil {
			return nil, err
		}
		result = append(result, struct {
			Name  string  `json:"name"`
			Color string  `json:"color"`
			Total float32 `json:"total"`
		}{
			Name:  name,
			Color: color,
			Total: total,
		})
	}
	return result, nil
}

func GetCategoryColor(name string) (string, error) {
	sqlScript := `
		SELECT color
		FROM categories	
		WHERE name = $1
	`
	var color string
	err := Pool.QueryRow(context.Background(), sqlScript, name).Scan(&color)
	if err != nil {
		return "", err
	}
	return color, nil
}

func GetExpensesByCategory(category_id int) ([]models.Expense, error) {
	sqlScript := `
		SELECT e.id, e.user_id, e.category_id, e.amount, e.name, e.description, e.expense_date, e.expense_time
		FROM expenses e
		WHERE category_id = $1
	`
	rows, err := Pool.Query(context.Background(), sqlScript, category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var expense models.Expense

		err := rows.Scan(&expense.Id, &expense.User_id, &expense.Category_id,
			&expense.Amount, &expense.Name, &expense.Description, &expense.Expense_date, &expense.Expense_time)
		if err != nil {
			return nil, err
		}
		result = append(result, expense)
	}
	return result, nil
}
