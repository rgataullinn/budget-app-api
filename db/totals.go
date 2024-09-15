package db

import "context"

func GetTotalSpentFromDb() (float32, error) {
	sqlScript := `
		SELECT COALESCE(SUM(amount), 0) as total_spent
		FROM expenses;
	`
	var total_spent float32
	err := Pool.QueryRow(context.Background(), sqlScript).Scan(&total_spent)
	if err != nil {
		return 0, err
	}
	return total_spent, nil
}
