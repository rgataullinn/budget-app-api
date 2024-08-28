package handlers

import (
	"net/http"
	"personal-finance-api/db"
	"personal-finance-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	if category.Id != 0 {
		err := db.UpdateCategoryInDb(category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
		return
	}
	err := db.AddCategoryInDb(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category added"})
}

func GetCategory(c *gin.Context) {
	idParam := c.Query("id")

	// Convert the 'id' to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	category, err := db.GetCategoryInDb(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category})
}

func DeleteCategory(c *gin.Context) {
	idParam := c.Query("id")

	// Convert the 'id' to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = db.DeleteCategoryInDb(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
