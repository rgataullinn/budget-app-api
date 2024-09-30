package models

import "time"

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Expense struct {
	Id           int       `json:"id"`
	User_id      int       `json:"user_id"`
	Category_id  int       `json:"category_id"`
	Category     string    `json:"category_name"`
	Amount       float32   `json:"amount"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Expense_date string    `json:"expense_date"`
	Expense_time string    `json:"expense_time"`
	Color        string    `json:"color"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
