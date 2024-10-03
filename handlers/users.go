package handlers

import (
	"net/http"
	"personal-finance-api/models"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"personal-finance-api/db"

	"github.com/gin-gonic/gin"
)

func isValidData(user models.User) (bool, string) {
	username := user.Username
	password := user.Password
	email := user.Email
	if len(username) < 3 || len(username) > 20 || strings.Contains(username, " ") || !unicode.IsLetter(rune(username[0])) {
		return false, "Username is not valid. It must be 3-20 characters long, start with a letter, and contain no spaces."
	}

	if len(password) < 6 || strings.Contains(password, " ") {
		return false, "Password is not valid. It must be at least 6 characters long and contain no spaces."
	}

	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return false, "Email is not valid. Please enter a valid email address."
	}

	return true, ""
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	isValid, errMsg := isValidData(newUser)

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	isExist, err := db.IsExist(newUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if isExist {
		c.JSON(http.StatusConflict, gin.H{"error": "User with the same username already exists"})
		return
	}

	err = db.AddUserInDb(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User successfully created"})
}

func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	isExist, err := db.IsExist(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if !isExist {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this username don't exist"})
		return
	}

	isValid, err := db.ValidateUserCredentials(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password:("})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func DeleteUser(c *gin.Context) {
	idParam := c.Query("id")

	// Convert the 'id' to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = db.DeleteUserInDb(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successful"})

}
