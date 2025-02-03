package main

import (
	"github.com/gin-gonic/gin"
)

func (app *Config) routes(router *gin.Engine) {
	// Настройка CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Link")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "300")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Middleware для проверки работоспособности
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Маршруты
	router.POST("/depositMoney", app.depositMoney)
	router.POST("/transferMoney", app.transferMoney)
	router.GET("/getLastTransactions", app.GetLastTransactions)
}
