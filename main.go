package main

import (
	"log"
	"os"
	"personal-finance-api/db"
	"personal-finance-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	router := gin.Default()
	router.POST("/users", handlers.CreateUser)
	router.GET("/users", handlers.LoginUser)
	router.DELETE("/users", handlers.DeleteUser)

	router.POST("/expenses", handlers.PostExpense)
	router.GET("/expenses", handlers.GetExpense)
	router.DELETE("/expenses", handlers.DeleteExpense)

	router.POST("/category", handlers.PostCategory)
	router.GET("/category", handlers.GetCategory)
	router.DELETE("/category", handlers.DeleteCategory)

	connString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName
	db.ConnectDB(connString)
	defer db.CloseDB()

	router.Run(":8080")
}
