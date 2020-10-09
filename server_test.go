package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

var server *Server
var createdAccount = AccountEventResponse{
	Destination: Account{
		ID:      100,
		Balance: 100,
	},
}

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DBTEST_URL")
	log.Println(dbURL)
	server = NewServer("EBANX Assignment")
	server.init()
	code := m.Run()
	os.Exit(code)
}

func TestReset(t *testing.T) {
	request := criarRequisicao("POST", "/reset", nil)
	response := executarRequisicao(request)
	assert.Equal(t, http.StatusOK, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "OK", string(body), "should be equal")
}

func TestGetBalanceNonExistingAccount(t *testing.T) {
	request := criarRequisicao("GET", "/balance?account_id=1234", nil)
	response := executarRequisicao(request)
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
	request := criarRequisicao("POST", "/event", bodyPost)
	response := executarRequisicao(request)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	var retorno AccountEventResponse
	json.Unmarshal(body, &retorno)
	assert.Equal(t, createdAccount, retorno, "should be equal")
}

func TestRetornarTodas(t *testing.T) {
	request := criarRequisicao("GET", "/", nil)
	response := executarRequisicao(request)
	assert.Equal(t, http.StatusOK, response.StatusCode, "shoud be equal")
	body, _ := ioutil.ReadAll(response.Body)
	var retorno map[string]string
	json.Unmarshal(body, &retorno)
	assert.Equal(t, "Retornando...", retorno["msg"], "should be equal")
}

func criarRequisicao(metodo string, url string, reqBody []byte) *http.Request {
	var req *http.Request
	if reqBody != nil {
		req, _ = http.NewRequest(metodo, url, strings.NewReader(string(reqBody)))
	} else {
		req, _ = http.NewRequest(metodo, url, nil)
	}
	return req
}

func executarRequisicao(req *http.Request) *http.Response {
	req.Header.Set("Content-Type", "application/json")
	ret, _ := server.Fiber.Test(req)
	return ret
}
