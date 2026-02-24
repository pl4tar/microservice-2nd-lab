package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

func ValidateToken(tokenString string) (bool, error) {
	authServiceURL := "http://localhost:8081/validate-token"
	data := map[string]string{"token": tokenString}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(authServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}
		valid, err := ValidateToken(tokenString)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetData(c *gin.Context) {
	data := []Data{
		{ID: 1, Content: "Protected Data 1"},
		{ID: 2, Content: "Protected Data 2"},
	}

	c.JSON(http.StatusOK, data)
}

func main() {
	r := gin.Default()

	r.Use(TokenAuthMiddleware()) // Защита всех маршрутов данным сервисом

	r.GET("/data", GetData)

	r.Run(":8082")
}
