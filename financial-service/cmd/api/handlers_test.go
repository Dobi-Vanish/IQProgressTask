package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"financial-service/data"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
	data.PostgresTestRepository
}

func (m *MockRepository) GetLastTransactions(id int) ([]*Transactions, error) {
	args := m.Called(id)
	return args.Get(0).([]*Transactions), args.Error(1)
}

func (m *MockRepository) AddMoney(id int, amount decimal.Decimal) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func (m *MockRepository) AddTransaction(amount decimal.Decimal, id ...int) error {
	args := m.Called(amount, id)
	return args.Error(0)
}

func (m *MockRepository) DecreaseMoney(id int, amount decimal.Decimal) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

// TestGetLastTransactions_InvalidJSON checks invalid JSON
func TestGetLastTransactions_InvalidJSON(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	invalidJSON := `{"Id": "invalid"}`

	req, _ := http.NewRequest("GET", "/getLastTransactions", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Invalid JSON", response["message"])
}

// TestGetLastTransactions_Success check success receiving
func TestGetLastTransactions_Success(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	transactions := []*Transactions{
		{ID: 1, UserIDSource: 123, UserIDEndpoint: 456, Amount: decimal.NewFromFloat(100.50), CreatedAt: time.Now()},
		{ID: 2, UserIDSource: 123, UserIDEndpoint: 789, Amount: decimal.NewFromFloat(200.75), CreatedAt: time.Now()},
	}

	mockRepo.On("GetLastTransactions", 123).Return(transactions, nil)

	payload := map[string]interface{}{
		"Id": 123,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("GET", "/getLastTransactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, false, response["error"])
	assert.Equal(t, "Fetched all transactions", response["message"])
	assert.Equal(t, transactions, response["data"])

	mockRepo.AssertCalled(t, "GetLastTransactions", 123)
}

// TestGetLastTransactions_Error checks error handling
func TestGetLastTransactions_Error(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	mockRepo.On("GetLastTransactions", 123).Return([]*Transactions{}, errors.New("database error"))

	payload := map[string]interface{}{
		"Id": 123,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("GET", "/getLastTransactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Couldn't fetch last 10 transactions ", response["message"])

	mockRepo.AssertCalled(t, "GetLastTransactions", 123)
}

// TestDepositMoney_InvalidJSON checks invalid JSON
func TestDepositMoney_InvalidJSON(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	invalidJSON := `{"Amount": "invalid", "Id": 123}`

	req, _ := http.NewRequest("POST", "/depositMoney", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Invalid JSON", response["message"])
}

// TestDepositMoney_Success check success deposit
func TestDepositMoney_Success(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	id := 123

	mockRepo.On("AddMoney", id, amount).Return(nil)

	mockRepo.On("AddTransaction", amount, []int{id}).Return(nil)

	payload := map[string]interface{}{
		"Amount": amount.String(),
		"Id":     id,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/depositMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, false, response["error"])
	assert.Equal(t, fmt.Sprintf("Deposit money worked for user with id %d, added money %s", id, amount.String()), response["message"])

	mockRepo.AssertCalled(t, "AddMoney", id, amount)
	mockRepo.AssertCalled(t, "AddTransaction", amount, []int{id})
}

// TestDepositMoney_AddMoneyError check error while depositing
func TestDepositMoney_AddMoneyError(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	id := 123

	mockRepo.On("AddMoney", id, amount).Return(errors.New("database error"))

	payload := map[string]interface{}{
		"Amount": amount.String(),
		"Id":     id,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/depositMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Couldn't add money to the user", response["message"])

	mockRepo.AssertCalled(t, "AddMoney", id, amount)
}

// TestDepositMoney_AddTransactionError error handling during depositing
func TestDepositMoney_AddTransactionError(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	id := 123

	mockRepo.On("AddMoney", id, amount).Return(nil)

	mockRepo.On("AddTransaction", amount, []int{id}).Return(errors.New("database error"))

	payload := map[string]interface{}{
		"Amount": amount.String(),
		"Id":     id,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/depositMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Couldn't add transaction", response["message"])

	mockRepo.AssertCalled(t, "AddMoney", id, amount)
	mockRepo.AssertCalled(t, "AddTransaction", amount, []int{id})
}

// TestTransferMoney_InvalidJSON checks invalid JSON
func TestTransferMoney_InvalidJSON(t *testing.T) {
	router := gin.Default()
	app := &Config{}
	app.routes(router)

	invalidJSON := `{"Amount": "invalid", "IdSource": 123, "IdEndpoint": 456}`

	req, _ := http.NewRequest("POST", "/transferMoney", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Invalid JSON", response["message"])
}

// TestTransferMoney_Success success money transfer
func TestTransferMoney_Success(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	idSource := 123
	idEndpoint := 456

	mockRepo.On("DecreaseMoney", idSource, amount).Return(nil)

	mockRepo.On("AddMoney", idEndpoint, amount).Return(nil)

	mockRepo.On("AddTransaction", amount, []int{idSource, idEndpoint}).Return(nil)

	payload := map[string]interface{}{
		"Amount":     amount.String(),
		"IdSource":   idSource,
		"IdEndpoint": idEndpoint,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/transferMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, false, response["error"])
	assert.Equal(t, "Transfer money worked successfully", response["message"])

	mockRepo.AssertCalled(t, "DecreaseMoney", idSource, amount)
	mockRepo.AssertCalled(t, "AddMoney", idEndpoint, amount)
	mockRepo.AssertCalled(t, "AddTransaction", amount, []int{idSource, idEndpoint})
}

// TestTransferMoney_DecreaseMoneyError error handling during decreasing money
func TestTransferMoney_DecreaseMoneyError(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	idSource := 123
	idEndpoint := 456

	mockRepo.On("DecreaseMoney", idSource, amount).Return(errors.New("database error"))

	payload := map[string]interface{}{
		"Amount":     amount.String(),
		"IdSource":   idSource,
		"IdEndpoint": idEndpoint,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/transferMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Couldn't decrease money from the source user", response["message"])

	mockRepo.AssertCalled(t, "DecreaseMoney", idSource, amount)
}

// TestTransferMoney_AddMoneyError error handling during adding money to balance
func TestTransferMoney_AddMoneyError(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	idSource := 123
	idEndpoint := 456

	mockRepo.On("DecreaseMoney", idSource, amount).Return(nil)

	mockRepo.On("AddMoney", idEndpoint, amount).Return(errors.New("database error"))

	payload := map[string]interface{}{
		"Amount":     amount.String(),
		"IdSource":   idSource,
		"IdEndpoint": idEndpoint,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/transferMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Couldn't add money to the destination user", response["message"])

	mockRepo.AssertCalled(t, "DecreaseMoney", idSource, amount)
	mockRepo.AssertCalled(t, "AddMoney", idEndpoint, amount)
}

// TestTransferMoney_AddTransactionError error handling during adding transaction
func TestTransferMoney_AddTransactionError(t *testing.T) {
	router := gin.Default()

	mockRepo := new(MockRepository)
	app := &Config{Repo: mockRepo}
	app.routes(router)

	amount := decimal.NewFromFloat(100.50)
	idSource := 123
	idEndpoint := 456

	mockRepo.On("DecreaseMoney", idSource, amount).Return(nil)

	mockRepo.On("AddMoney", idEndpoint, amount).Return(nil)

	mockRepo.On("AddTransaction", amount, []int{idSource, idEndpoint}).Return(errors.New("database error"))

	payload := map[string]interface{}{
		"Amount":     amount.String(),
		"IdSource":   idSource,
		"IdEndpoint": idEndpoint,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/transferMoney", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, true, response["error"])
	assert.Equal(t, "Couldn't add transaction", response["message"])

	mockRepo.AssertCalled(t, "DecreaseMoney", idSource, amount)
	mockRepo.AssertCalled(t, "AddMoney", idEndpoint, amount)
	mockRepo.AssertCalled(t, "AddTransaction", amount, []int{idSource, idEndpoint})
}
