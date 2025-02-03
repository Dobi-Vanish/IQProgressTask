package main

import (
	"bytes"
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

type testRoundTripper struct {
	fn RoundTripFunc
}

type MockRepo struct {
	mock.Mock
}

func (t testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req), nil
}
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: testRoundTripper{fn: fn},
	}
}

func (m *MockRepo) GetLastTransactions(id int) ([]Transaction, error) {
	args := m.Called(id)
	return args.Get(0).([]Transaction), args.Error(1)
}

func TestGetLastTransactions_Success(t *testing.T) {
	// Настройка моков
	repo := new(MockRepo)
	app := &Config{Repo: repo}

	// Ожидаемые данные
	expectedTransactions := []Transaction{
		{ID: 1, Amount: decimal.NewFromInt(100)},
		{ID: 2, Amount: decimal.NewFromInt(200)},
	}
	repo.On("GetLastTransactions", 1).Return(expectedTransactions, nil)

	// Создание запроса
	requestBody := map[string]interface{}{"Id": 1}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("GET", "/transactions", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Вызов хендлера
	app.GetLastTransactions(rr, req)

	// Проверка результата
	assert.Equal(t, http.StatusAccepted, rr.Code)

	var response jsonResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.False(t, response.Error)
	assert.Equal(t, "Fetched all transactions", response.Message)
	assert.Equal(t, expectedTransactions, response.Data)

	// Проверка, что мок был вызван
	repo.AssertCalled(t, "GetLastTransactions", 1)
}

func TestDepositMoney_Success(t *testing.T) {
	// Настройка моков
	repo := new(MockRepo)
	app := &Config{Repo: repo}

	// Ожидаемые данные
	repo.On("AddMoney", 1, decimal.NewFromInt(100)).Return(nil)

	// Создание запроса
	requestBody := map[string]interface{}{"Id": 1, "Amount": 100}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Вызов хендлера
	app.depositMoney(rr, req)

	// Проверка результата
	assert.Equal(t, http.StatusAccepted, rr.Code)

	var response jsonResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.False(t, response.Error)
	assert.Contains(t, response.Message, "deposit money worked for user with id 1")

	// Проверка, что мок был вызван
	repo.AssertCalled(t, "AddMoney", 1, decimal.NewFromInt(100))
}

func TestTransferMoney_Success(t *testing.T) {
	// Настройка моков
	repo := new(MockRepo)
	app := &Config{Repo: repo}

	// Ожидаемые данные
	repo.On("AddMoney", 2, decimal.NewFromInt(100)).Return(nil)
	repo.On("DecreaseMoney", 1, decimal.NewFromInt(100)).Return(nil)

	// Создание запроса
	requestBody := map[string]interface{}{"IdSource": 1, "IdEndpoint": 2, "Amount": 100}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Вызов хендлера
	app.transferMoney(rr, req)

	// Проверка результата
	assert.Equal(t, http.StatusAccepted, rr.Code)

	var response jsonResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.False(t, response.Error)
	assert.Contains(t, response.Message, "transfer money worked successfully")

	// Проверка, что моки были вызваны
	repo.AssertCalled(t, "AddMoney", 2, decimal.NewFromInt(100))
	repo.AssertCalled(t, "DecreaseMoney", 1, decimal.NewFromInt(100))
}
