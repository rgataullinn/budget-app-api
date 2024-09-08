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

func GetAllExpense(c *gin.Context) {
	isCategory, err1 := strconv.ParseBool(c.Query("isCategory"))
	isDate, err2 := strconv.ParseBool(c.Query("isDate"))

	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1})
		return
	}

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2})
		return
	}

	if isCategory && isDate {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "choose only one type"})
		return
	}

	if isCategory {
		expenses, err := db.GetAllExpenseByCategory()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"expenses": expenses})
		return
	}

	if isDate {
		expenses, err := db.GetAllExpensesByDate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"expenses": expenses})
		return
	}
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
