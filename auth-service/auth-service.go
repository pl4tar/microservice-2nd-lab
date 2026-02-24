package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "log"
	"net/http"
	"time"
)

var jwtKey = []byte("my_secret_key")

// Структура для входных данных пользователя
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Структура для формирования JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Функция для генерации JWT
func GenerateJWT(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Простой проверочный механизм для пользователя (должен быть заменен на реальную проверку)
	if creds.Username != "admin" || creds.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Новый эндпоинт для проверки токена
func ValidateToken(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tokenString := req["token"]
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}

func main() {
	r := gin.Default()

	r.POST("/login", GenerateJWT)
	r.POST("/validate-token", ValidateToken) // Новый эндпоинт для проверки токена

	r.Run(":8081")
}
