package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server *Server

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DBTEST_URL")
	log.Println(dbURL)
	server = NewServer("Cardea Servico MS")
	server.init()
	code := m.Run()
	os.Exit(code)
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
