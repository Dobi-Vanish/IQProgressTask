package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCORS checks for CORS middleware correctly set headers
func TestCORS(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	req, _ := http.NewRequest("OPTIONS", "/ping", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Accept, Authorization, Content-Type, X-CSRF-Token", resp.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", resp.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "300", resp.Header().Get("Access-Control-Max-Age"))

	assert.Equal(t, 204, resp.Code)
}

// TestPingRoute  checks routes for correct handling
func TestPingRoute(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	req, _ := http.NewRequest("GET", "/ping", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "pong", response["message"])
}

// TestDepositMoneyRoute checks routes for correct handling
func TestDepositMoneyRoute(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	payload := map[string]interface{}{
		"Id":     123,
		"Amount": "100.50",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/depositMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

// TestTransferMoneyRoute checks routes for correct handling
func TestTransferMoneyRoute(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	payload := map[string]interface{}{
		"IdSource":   123,
		"IdEndpoint": 456,
		"Amount":     "50.25",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/transferMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

// TestGetLastTransactionsRoute checks routes for correct handling
func TestGetLastTransactionsRoute(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	req, _ := http.NewRequest("GET", "/getLastTransactions?id=123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}
