package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

// GetLastTransactions retrieves all users from the database, sort them by points
func (app *Config) GetLastTransactions(c *gin.Context) {
	// Структура для парсинга JSON-запроса
	var requestPayload struct {
		ID int `json:"Id"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		// Возвращаем ошибку, если JSON невалидный
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	// Получаем транзакции из репозитория
	transactions, err := app.Repo.GetLastTransactions(requestPayload.ID)
	if err != nil {
		// Возвращаем ошибку, если не удалось получить транзакции
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't fetch last 10 transactions",
		})
		return
	}

	// Возвращаем успешный ответ с данными
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Fetched all transactions",
		"data":    transactions,
	})
}

func (app *Config) depositMoney(c *gin.Context) {
	// Структура для парсинга JSON-запроса
	var requestPayload struct {
		Amount decimal.Decimal `json:"Amount"`
		ID     int             `json:"Id"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		// Возвращаем ошибку, если JSON невалидный
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	// Добавляем деньги пользователю
	err := app.Repo.AddMoney(requestPayload.ID, requestPayload.Amount)
	if err != nil {
		// Возвращаем ошибку, если не удалось добавить деньги
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't add money to the user",
		})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": fmt.Sprintf("Deposit money worked for user with id %d, added money %s", requestPayload.ID, requestPayload.Amount.String()),
	})
}

func (app *Config) transferMoney(c *gin.Context) {
	// Структура для парсинга JSON-запроса
	var requestPayload struct {
		Amount     decimal.Decimal `json:"Amount"`
		IDSource   int             `json:"IdSource"`
		IDEndpoint int             `json:"IdEndpoint"`
	}

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		// Возвращаем ошибку, если JSON невалидный
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	// Переводим деньги от одного пользователя к другому
	err := app.Repo.AddMoney(requestPayload.IDEndpoint, requestPayload.Amount)
	if err != nil {
		// Возвращаем ошибку, если не удалось добавить деньги
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't add money to the destination user",
		})
		return
	}

	err = app.Repo.DecreaseMoney(requestPayload.IDSource, requestPayload.Amount)
	if err != nil {
		// Возвращаем ошибку, если не удалось списать деньги
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't decrease money from the source user",
		})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Transfer money worked successfully",
	})
}
