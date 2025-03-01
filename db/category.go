package db

import (
	"context"
	"personal-finance-api/models"
)

func CreateCategory(newCategory models.Category) error {
	sqlScript := `
		INSERT INTO categories (name, description)
		VALUES($1, $2, $3);
	`
	_, err := Pool.Exec(context.Background(), sqlScript,
		newCategory.Name,
		newCategory.Description)

	return err
}

func UpdateCategory(category models.Category) error {
	sqlScript := `
		UPDATE categories
		SET
			name = $2,
			description = $3,
		WHERE
			id = $1
`
	_, err := Pool.Exec(context.Background(), sqlScript,
		category.Id,
		category.Name,
		category.Description)

	return err
}

func GetCategory(id int) (models.Category, error) {
	sqlScript := `
		SELECT id, name, description
		FROM categories
		WHERE id= $1
	`
	var category models.Category
	row := Pool.QueryRow(context.Background(), sqlScript, id)
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func DeleteCategory(id int) error {
	sqlScript := `
		DELETE FROM categories
		WHERE id = $1
	`
	_, err := Pool.Exec(context.Background(), sqlScript, id)
	return err
}

func GetCategories(month int) ([]models.Category, error) {
	sqlScript := `
		SELECT c.id, c.name, c.description, sum(e.amount)
		FROM categories c
		JOIN expenses e
		ON c.id = e.category_id
		WHERE EXTRACT(MONTH FROM e.expense_date::date) = $1
		GROUP BY c.id
	`
	rows, err := Pool.Query(context.Background(), sqlScript, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.Total)
		if err != nil {
			return nil, err
		}
		result = append(result, category)
	}
	return result, nil
}
