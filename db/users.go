package db

import (
	"context"
	"personal-finance-api/models"
)

func AddUserInDb(newUser models.User) error {
	sqlScript := `
        INSERT INTO users (username, password, email) 
        VALUES ($1, $2, $3);
    `
	_, err := Pool.Exec(context.Background(), sqlScript,
		newUser.Username, newUser.Password, newUser.Email)
	if err != nil {
		return err
	}
	return nil
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

func GetUserPassword(username string) (int, string, error) {
	sqlScript := `
		select id, password
		from users
		where username = $1
	`
	var id int
	var password string
	err := Pool.QueryRow(context.Background(), sqlScript, username).Scan(&id, &password)
	if err != nil {
		return -1, "", err
	}
	return id, password, nil
}

func GetUser(id int) (models.User, error) {
	sqlScript := `
		select id, username, email
		from users
		where id = $1
	`
	var user models.User
	err := Pool.QueryRow((context.Background()),
		sqlScript, id).Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
