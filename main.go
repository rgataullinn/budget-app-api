package main

import (
	"os"
	"personal-finance-api/db"
	"personal-finance-api/handlers"
	"personal-finance-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/users", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.DELETE("/users", handlers.DeleteUser)
	router.GET("/validate", middleware.RequireAuth, handlers.Validate)

	router.POST("/expense", handlers.CreateExpense)
	router.PUT("/expense", handlers.UpdateExpense)
	router.GET("/expense", handlers.GetExpense)
	router.GET("/expensesByDate", handlers.GetAllDatesWithExpenses)
	router.GET("/expensesByCategory", handlers.GetAllCategoriesWithExpenses)
	router.DELETE("/expense", handlers.DeleteExpense)

	router.POST("/category", handlers.CreateCategory)
	router.PUT("/category", handlers.UpdateCategory)
	router.GET("/category", handlers.GetCategory)
	router.GET("/categories", handlers.GetCategories)
	router.GET("/categoriesList", handlers.GetCategoriesList)
	router.DELETE("/category", handlers.DeleteCategory)

	router.GET("/totalSpent", handlers.GetTotalSpent)
	// router.GET("/totatlEarn", handlers.GetTotalEarns)

	connString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	db.ConnectDB(connString)
	defer db.CloseDB()

	router.Run("0.0.0.0:8080")
}
