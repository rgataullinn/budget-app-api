package db

import (
	"context"
	"personal-finance-api/models"
)

func AddUserInDb(newUser models.User) error {
	sqlScript := `
        INSERT INTO users (username, password, email) 
        VALUES ($1, crypt($2, gen_salt('bf')), $3);
    `
	_, err := Pool.Exec(context.Background(), sqlScript,
		newUser.Username, newUser.Password, newUser.Email)
	if err != nil {
		return err
	}
	return nil
}

func Login(username string, password string) (bool, int, error) {
	sqlScript := `
        SELECT EXISTS (
            SELECT 1 FROM users WHERE username = $1 AND password = crypt($2, password)
        );
    `
	var exists bool
	err := Pool.QueryRow(context.Background(), sqlScript, username, password).Scan(&exists)
	if err != nil {
		return false, -1, err
	}

	if exists {
		id, err := getUserId(username)
		if err != nil {
			return false, -1, err
		}
		return true, id, nil
	}

	return exists, -1, nil
}

func getUserId(username string) (int, error) {
	sqlScript := `
	select id
	from users
	where username = $1;
	`

	var id int
	err := Pool.QueryRow(context.Background(), sqlScript, username).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func DeleteUserInDb(id int) error {
	sqlScript := `
		DELETE FROM users
		WHERE id = $1;
    `
	_, err := Pool.Exec(context.Background(), sqlScript, id)
	if err != nil {
		return err
	}
	return nil
}
