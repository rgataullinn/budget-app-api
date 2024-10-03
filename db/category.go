package db

import (
	"context"
	"personal-finance-api/models"
)

func CreateCategory(newCategory models.Category) error {
	sqlScript := `
		INSERT INTO categories (name, description, color)
		VALUES($1, $2, $3);
	`
	_, err := Pool.Exec(context.Background(), sqlScript,
		newCategory.Name,
		newCategory.Description,
		newCategory.Color)

	return err
}

func UpdateCategory(category models.Category) error {
	sqlScript := `
		UPDATE categories
		SET
			name = $2,
			description = $3,
			color = $4
		WHERE
			id = $1
`
	_, err := Pool.Exec(context.Background(), sqlScript,
		category.Id,
		category.Name,
		category.Description,
		category.Color)

	return err
}

func GetCategory(id int) (models.Category, error) {
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
		select c.name, sum(e.amount) as total, c.color
		from categories c 
		join expenses e 
		on c.id = e.category_id 
		WHERE EXTRACT(MONTH FROM e.expense_date::date) = $1
		group by c.id;
	`
	var result []models.Category
	rows, err := Pool.Query(context.Background(), sqlScript, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Name, &category.Total, &category.Color)
		if err != nil {
			return nil, err
		}
		result = append(result, category)
	}
	return result, nil
}

func GetCategoriesList() ([]string, error) {
	sqlScript := `
		SELECT name
		FROM categories
	`
	rows, err := Pool.Query(context.Background(), sqlScript)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var category_name string
		err := rows.Scan(&category_name)
		if err != nil {
			return nil, err
		}
		result = append(result, category_name)
	}
	return result, nil
}
