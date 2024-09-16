package handlers

import (
	"net/http"
	"personal-finance-api/db"
	"personal-finance-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	if expense.Id != 0 {
		err := db.UpdateExpenseInDb(expense)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Expense updated"})
		return
	}

	err := db.AddExpenseInDb(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Expense added"})

}

func GetExpense(c *gin.Context) {
	idParam := c.Query("id")

	// Convert the 'id' to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	expense, err := db.GetExpenseFromDb(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": expense})

}

func GetAllExpensesByDate(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	expenses, err := db.GetAllExpensesByDate(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func GetAllExpensesByCategory(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	expenses, err := db.GetAllExpenseByCategory(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func DeleteExpense(c *gin.Context) {
	idParam := c.Query("id")

	// Convert the 'id' to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = db.DeleteExpenseFromDb(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
