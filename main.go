package main

import (
	"log"
	"os"
	"personal-finance-api/db"
	"personal-finance-api/handlers"
	"personal-finance-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		c.Header("ExposeHeaders", "Set-Cookie")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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

	router.Use(CORSMiddleware())


	router.POST("/users", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.DELETE("/users", handlers.DeleteUser)
	router.GET("/validate", middleware.RequireAuth, handlers.Validate)

	router.POST("/expense", middleware.RequireAuth, handlers.CreateExpense)
	router.PUT("/expense", middleware.RequireAuth, handlers.UpdateExpense)
	router.GET("/expense", middleware.RequireAuth, handlers.GetExpense)
	router.GET("/expenses", middleware.RequireAuth, handlers.GetAllExpenses)
	router.DELETE("/expense", middleware.RequireAuth, handlers.DeleteExpense)

	router.POST("/category", handlers.CreateCategory)
	router.PUT("/category", handlers.UpdateCategory)
	router.GET("/category", handlers.GetCategory)
	router.GET("/categories", handlers.GetCategories)
	router.DELETE("/category", handlers.DeleteCategory)

	router.GET("/totalSpent", handlers.GetTotalSpent)
	// router.GET("/totatlEarn", handlers.GetTotalEarns)

	connString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	db.ConnectDB(connString)
	defer db.CloseDB()

	router.Run("0.0.0.0:8080")
}
