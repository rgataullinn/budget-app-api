package db

import "context"

func GetTotalSpentFromDb() (int, error) {
	sqlScript := `
		SELECT sum(amount) as total_spent
		FROM expenses;
	`
	var total_spent int
	err := Pool.QueryRow(context.Background(), sqlScript).Scan(&total_spent)
	if err != nil {
		return 0, err
	}
	return total_spent, nil
}
