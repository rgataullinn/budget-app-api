package db

import "context"

func GetTotalSpentFromDb(month int) (float32, error) {
	sqlScript := `
		SELECT COALESCE(SUM(amount), 0) as total_spent
		FROM expenses
		WHERE EXTRACT(MONTH FROM expense_date::date) = $1;
	`
	var total_spent float32
	err := Pool.QueryRow(context.Background(), sqlScript, month).Scan(&total_spent)
	if err != nil {
		return 0, err
	}
	return total_spent, nil
}
