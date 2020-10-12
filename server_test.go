package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var server *Server
var createdAccount = AccountEventResponse{
	Destination: &Account{
		ID:      "100",
		Balance: 10,
	},
}
var depositCreatedAccount = AccountEventResponse{
	Destination: &Account{
		ID:      "100",
		Balance: 20,
	},
}

var withdrawCreatedAccount = AccountEventResponse{
	Origin: &Account{
		ID:      "100",
		Balance: 15,
	},
}

var transferExistingAccount = AccountEventResponse{
	Origin: &Account{
		ID:      "100",
		Balance: 0,
	},
	Destination: &Account{
		ID:      "300",
		Balance: 15,
	},
}

func TestMain(m *testing.M) {
	server = NewServer("EBANX Assignment")
	server.init()
	code := m.Run()
	os.Exit(code)
}

func TestReset(t *testing.T) {
	request := createRequest("POST", "/reset", nil)
	response := execRequest(request)
	assert.Equal(t, http.StatusOK, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "OK", string(body), "should be equal")
}

func TestGetBalanceNonExistingAccount(t *testing.T) {
	request := createRequest("GET", "/balance?account_id=1234", nil)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "0", string(body), "should be equal")
}

func TestCreateAccountWithInitialBalance(t *testing.T) {
	event := AccountEvent{
		Type:        "deposit",
		Destination: "100",
		Amount:      10}
	bodyPost, _ := json.Marshal(event)
	request := createRequest("POST", "/event", bodyPost)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	var retorno AccountEventResponse
	json.Unmarshal(body, &retorno)
	assert.Equal(t, createdAccount, retorno, "should be equal")
}

func TestDepositIntoExistingAccount(t *testing.T) {
	event := AccountEvent{
		Type:        "deposit",
		Destination: "100",
		Amount:      10}
	bodyPost, _ := json.Marshal(event)
	request := createRequest("POST", "/event", bodyPost)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	var retorno AccountEventResponse
	json.Unmarshal(body, &retorno)
	assert.Equal(t, depositCreatedAccount, retorno, "should be equal")
}

func TestGetBalanceExistingAccount(t *testing.T) {
	request := createRequest("GET", "/balance?account_id=100", nil)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "20", string(body), "should be equal")
}

func TestWithdrawNonExistingAccount(t *testing.T) {
	event := AccountEvent{
		Type:   "withdraw",
		Origin: "200",
		Amount: 10}
	bodyPost, _ := json.Marshal(event)
	request := createRequest("POST", "/event", bodyPost)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "0", string(body), "should be equal")
}

func TestWithdrawExistingAccount(t *testing.T) {
	event := AccountEvent{
		Type:   "withdraw",
		Origin: "100",
		Amount: 5}
	bodyPost, _ := json.Marshal(event)
	request := createRequest("POST", "/event", bodyPost)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	var retorno AccountEventResponse
	json.Unmarshal(body, &retorno)
	assert.Equal(t, withdrawCreatedAccount, retorno, "should be equal")
}

func TestTransferExistingAccount(t *testing.T) {
	event := AccountEvent{
		Type:        "transfer",
		Origin:      "100",
		Destination: "300",
		Amount:      15}
	bodyPost, _ := json.Marshal(event)
	request := createRequest("POST", "/event", bodyPost)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	var retorno AccountEventResponse
	json.Unmarshal(body, &retorno)
	assert.Equal(t, transferExistingAccount, retorno, "should be equal")
}

func TestTransferFromNonExistingAccount(t *testing.T) {
	event := AccountEvent{
		Type:        "transfer",
		Origin:      "200",
		Destination: "300",
		Amount:      15}
	bodyPost, _ := json.Marshal(event)
	request := createRequest("POST", "/event", bodyPost)
	response := execRequest(request)
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "0", string(body), "should be equal")
}

func createRequest(metodo string, url string, reqBody []byte) *http.Request {
	var req *http.Request
	if reqBody != nil {
		req, _ = http.NewRequest(metodo, url, strings.NewReader(string(reqBody)))
	} else {
		req, _ = http.NewRequest(metodo, url, nil)
	}
	return req
}

func execRequest(req *http.Request) *http.Response {
	req.Header.Set("Content-Type", "application/json")
	ret, _ := server.Fiber.Test(req)
	return ret
}
