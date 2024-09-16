package handlers

import (
	"net/http"
	"personal-finance-api/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTotalSpent(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	res, err := db.GetTotalSpentFromDb(month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"totalSpent": res})
}
