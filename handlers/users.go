package handlers

import (
	"net/http"
	"os"
	"personal-finance-api/models"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"personal-finance-api/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashed_password)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read username and password"})
		return
	}

	userId, hashedPassword, err := db.GetUserPassword(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password or username"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password or username"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password or username"})
		return
	}

	// secureFlag := true
	// if strings.HasPrefix(c.Request.Host, "0.0.0.0") || strings.HasPrefix(c.Request.Host, "localhost") {
	// 	secureFlag = false
	// }

	c.SetCookie(
		"Auth",
		tokenString,
		3600*2,
		"",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"id":      userId,
		"message": "Login successful",
	})
}

func DeleteUser(c *gin.Context) {
	idParam := c.Query("id")

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

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"userId": user})
}
