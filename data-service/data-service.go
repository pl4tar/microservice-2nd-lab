package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

// Модель данных для таблицы
type Data struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

var db *gorm.DB

// Функция для инициализации базы данных
func initDB() {
	var err error
	// Подключение к базе данных PostgreSQL
	dsn := "host=localhost user=postgres password=8967451230 dbname=data_service port=5432 sslmode=disable"
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	fmt.Println("Successfully connected to the database!")

	// Миграция: создаем таблицы в базе данных на основе моделей
	db.AutoMigrate(&Data{})
}

// Функция для получения всех данных из базы
func GetData(c *gin.Context) {
	var data []Data
	if err := db.Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// Функция для добавления данных в базу
func AddData(c *gin.Context) {
	var data Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Добавляем данные в базу
	if err := db.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func main() {
	// Инициализация базы данных
	initDB()

	// Настройка Gin для обработки запросов
	r := gin.Default()

	// Роуты для работы с данными
	r.GET("/data", GetData)  // Получить все данные
	r.POST("/data", AddData) // Добавить данные

	// Запуск сервера
	r.Run(":8082")
}
