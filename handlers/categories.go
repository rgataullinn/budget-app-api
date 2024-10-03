package handlers

import (
	"net/http"
	"personal-finance-api/db"
	"personal-finance-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error()})
		return
	}

	err := db.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create category in db",
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category added"})
}

func UpdateCategory(c *gin.Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error()})
		return
	}
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error()})
		return
	}
	category.Id = id

	err = db.UpdateCategory(category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update category in db",
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

func GetCategory(c *gin.Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error()})
		return
	}

	category, err := db.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get category from db",
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category})
}

func DeleteCategory(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"details": err.Error()})
		return
	}

	err = db.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete category from db",
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}

func GetCategories(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "wrong month format",
			"details": err.Error()})
		return
	}
	categories, err := db.GetCategories(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get categories from db",
			"details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
