package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type Transactions struct {
	ID             int             `json:"ID"`
	UserIDSource   int             `json:"UserIDSource"`
	UserIDEndpoint int             `json:"UserIDEndpoint"`
	Amount         decimal.Decimal `json:"Amount"`
	CreatedAt      time.Time       `json:"CreatedAt"`
}

// GetLastTransactions retrieves 10 last transactions for user from the database, sort them by points
func (app *Config) GetLastTransactions(c *gin.Context) {
	var requestPayload struct {
		ID int `json:"Id"`
	}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	transactions, err := app.Repo.GetLastTransactions(requestPayload.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't fetch last 10 transactions ",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Fetched all transactions",
		"data":    transactions,
	})
}

// depositMoney deposits money to the users balance
func (app *Config) depositMoney(c *gin.Context) {
	var requestPayload struct {
		Amount decimal.Decimal `json:"Amount"`
		ID     int             `json:"Id"`
	}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	err := app.Repo.AddMoney(requestPayload.ID, requestPayload.Amount)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't add money to the user",
		})
		return
	}
	err = app.Repo.AddTransaction(requestPayload.Amount, requestPayload.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't add transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": fmt.Sprintf("Deposit money worked for user with id %d, added money %s", requestPayload.ID, requestPayload.Amount.String()),
	})
}

// transferMoney transfers money from one user to another
func (app *Config) transferMoney(c *gin.Context) {
	var requestPayload struct {
		Amount     decimal.Decimal `json:"Amount"`
		IDSource   int             `json:"IdSource"`
		IDEndpoint int             `json:"IdEndpoint"`
	}

	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}
	err := app.Repo.DecreaseMoney(requestPayload.IDSource, requestPayload.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't decrease money from the source user",
		})
		return
	}

	err = app.Repo.AddMoney(requestPayload.IDEndpoint, requestPayload.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't add money to the destination user",
		})
		return
	}

	err = app.Repo.AddTransaction(requestPayload.Amount, requestPayload.IDSource, requestPayload.IDEndpoint)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Couldn't add transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Transfer money worked successfully",
	})
}
