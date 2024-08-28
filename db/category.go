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

func DeleteCategoryInDb(id int) error {
	sqlScript := `
		DELETE FROM categories
		WHERE id = $1
	`
	_, err := Pool.Exec(context.Background(), sqlScript, id)
	return err
}
