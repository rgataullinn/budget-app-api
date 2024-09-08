package main

import (
	"log"
	"os"
	"personal-finance-api/db"
	"personal-finance-api/handlers"

	"github.com/gin-contrib/cors"
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
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/users", handlers.CreateUser)
	router.GET("/users", handlers.LoginUser)
	router.DELETE("/users", handlers.DeleteUser)

	router.POST("/expense", handlers.PostExpense)
	router.GET("/expense", handlers.GetExpense)
	router.GET("/expenses", handlers.GetAllExpense)
	router.DELETE("/expense", handlers.DeleteExpense)

	router.POST("/category", handlers.PostCategory)
	router.GET("/category", handlers.GetCategory)
	router.DELETE("/category", handlers.DeleteCategory)
	router.GET("/categories", handlers.GetAllCategories)

	connString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	db.ConnectDB(connString)
	defer db.CloseDB()

	router.Run("0.0.0.0:8080")
}
