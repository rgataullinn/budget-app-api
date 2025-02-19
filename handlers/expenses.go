package handlers

import (
	"net/http"
	"personal-finance-api/db"
	"personal-finance-api/helpers"
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

	userId, err := helpers.GetUserIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to get userId from cookie",
			"detail": err.Error(),
		})
		return
	}
	expense.User_id = userId
	err = db.CreateExpense(expense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to add expense into db",
			"detail": err.Error()})
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

	userId, err := helpers.GetUserIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to get userId from cookie",
			"detail": err.Error(),
		})
		return
	}
	expense.User_id = userId

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

func GetAllExpenses(c *gin.Context) {
	groupBy := c.Query("groupBy")

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "wrong month format",
			"details": err.Error()})
		return
	}

	user_id, err := helpers.GetUserIdFromCookie(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed to get userId from cookie",
			"detail": err.Error(),
		})
		return
	}

	switch groupBy {
	case "day":
		expensesGroupedByDay, err := db.GetAllExpensesGroupedByDay(month, user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to get expenses from db",
				"detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": expensesGroupedByDay})
	case "category":
		expensesGroupedByCategory, err := db.GetAllExpensesGroupedByCategory(month, user_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to get expenses from db",
				"detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": expensesGroupedByCategory})
	}
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
