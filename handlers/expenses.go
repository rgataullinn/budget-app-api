package handlers

import (
	"net/http"
	"personal-finance-api/db"
	"personal-finance-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error()})
		return
	}

	err := db.CreateExpense(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to add expense into db",
			"detail": err.Error})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Expense added"})
}

func UpdateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error()})
		return
	}
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error()})
		return
	}
	expense.Id = id

	err = db.UpdateExpense(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed update expense in db",
			"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Expense updated"})
}

func GetExpense(c *gin.Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error()})
		return
	}

	expense, err := db.GetExpense(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to get expense from db",
			"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": expense})
}

func GetAllCategoriesWithExpenses(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "wrong month format",
			"details": err.Error()})
		return
	}
	expenses, err := db.GetAllCategoriesWithExpenses(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to get expenses from db",
			"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func GetAllDatesWithExpenses(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "wrong month format",
			"detail": err.Error()})
		return
	}
	expenses, err := db.GetAllDatesWithExpenses(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to get expenses from db",
			"detail": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func DeleteExpense(c *gin.Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error()})
		return
	}

	err = db.DeleteExpenseFromDb(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete expense from db",
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
}
