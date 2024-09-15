package handlers

import (
	"net/http"
	"personal-finance-api/db"

	"github.com/gin-gonic/gin"
)

func GetTotalSpent(c *gin.Context) {
	res, err := db.GetTotalSpentFromDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"totalSpent": res})
}
