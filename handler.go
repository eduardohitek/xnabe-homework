package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Handler represents a Handler Object
type Handler struct {
	DB  *DB
	Log *LogHeimdall
}

// NewHandler returns a new Handler
func NewHandler(DB *DB, Log *LogHeimdall) *Handler {
	return &Handler{DB: DB, Log: Log}
}

func (h *Handler) handlerRetornarTodas(c *fiber.Ctx) error {
	c.JSON(fiber.Map{"msg": "Retornando..."})
	return nil
}

func (h *Handler) reset(c *fiber.Ctx) error {
	_, error := h.DB.resetDatabase()
	if error != nil {
		return error
	}
	c.SendString("OK")
	return nil
}

func (h *Handler) getBalance(c *fiber.Ctx) error {
	ID := c.Query("account_id", "0")
	log.Println("ID", ID)
	account, erro := h.DB.getBalance(ID)
	if erro != nil {
		c.Status(404).SendString("0")
		return nil
	}
	c.JSON(account.Balance)
	return nil
}

func (h *Handler) handleAccountEvent(c *fiber.Ctx) error {
	accountEvent, _ := recuperarDadosAccountEvent(c)
	switch accountEvent.Type {
	case "deposit":
		account, error := h.DB.getBalance(accountEvent.Destination)
		log.Println(account)
		if error != nil {
			newAccount, _ := createAccount(accountEvent, h.DB)
			log.Println("aqui", newAccount)
			c.Status(fiber.StatusCreated).JSON(AccountEventResponse{Destination: *newAccount})
		} else {
			newAccount, _ := depositAmount(accountEvent, h.DB)
			c.Status(fiber.StatusCreated).JSON(AccountEventResponse{Destination: *newAccount})
		}
	}
	return nil
}

func createAccount(event AccountEvent, DB *DB) (*Account, error) {
	account := Account{ID: event.Destination, Balance: event.Amount}
	_, error := DB.createAccount(account)
	newAccount, _ := DB.getBalance(account.ID)
	return newAccount, error
}

func depositAmount(event AccountEvent, DB *DB) (*Account, error) {
	account := Account{ID: event.Destination, Balance: event.Amount}
	error := DB.depositAmount(account)
	newAccount, _ := DB.getBalance(account.ID)
	return newAccount, error
}

// Error Return an Error response to the user
func Error(c *fiber.Ctx, msg string, erro string) {
	c.Status(500).JSON(ErrorResponse{Msg: msg, Erro: erro})
	return
}

func recuperarDadosAccountEvent(c *fiber.Ctx) (AccountEvent, error) {
	accountEvent := new(AccountEvent)
	erro := c.BodyParser(accountEvent)
	return *accountEvent, erro
}
