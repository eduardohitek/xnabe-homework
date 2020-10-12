package main

import (
	"encoding/json"

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
		depositAccount := handleDeposit(h.DB, accountEvent)
		c.Status(fiber.StatusCreated).SendString(jsonConvert(AccountEventResponse{Destination: depositAccount}))
	case "withdraw":
		withdrawAccount, error := handleWithdraw(h.DB, accountEvent)
		if error != nil {
			c.Status(404).SendString("0")
		} else {
			c.Status(fiber.StatusCreated).JSON(AccountEventResponse{Origin: withdrawAccount})
		}
	case "transfer":
		withdrawAccount, depositAccount, error := handleTransfer(h.DB, accountEvent)
		if error != nil {
			c.Status(404).SendString("0")
		} else {
			c.Status(fiber.StatusCreated).JSON(AccountEventResponse{Origin: withdrawAccount, Destination: depositAccount})
		}
	}
	return nil
}

func handleDeposit(DB *DB, accountEvent AccountEvent) *Account {
	_, error := DB.getBalance(accountEvent.Destination)
	var newAccount *Account
	if error != nil {
		newAccount, _ = createAccount(accountEvent, DB)
	} else {
		newAccount, _ = depositAmount(accountEvent, DB)
	}
	return newAccount
}

func handleWithdraw(DB *DB, accountEvent AccountEvent) (*Account, error) {
	_, error := DB.getBalance(accountEvent.Origin)
	if error != nil {
		return nil, error
	}
	newAccount, _ := withdrawAmount(accountEvent, DB)
	return newAccount, nil
}

func handleTransfer(DB *DB, accountEvent AccountEvent) (*Account, *Account, error) {
	_, error := DB.getBalance(accountEvent.Origin)
	if error != nil {
		return nil, nil, error
	}
	withdrawAccount, _ := handleWithdraw(DB, accountEvent)
	depositAccount := handleDeposit(DB, accountEvent)
	return withdrawAccount, depositAccount, nil
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

func withdrawAmount(event AccountEvent, DB *DB) (*Account, error) {
	account := Account{ID: event.Origin, Balance: event.Amount}
	error := DB.withdrawAmount(account)
	newAccount, _ := DB.getBalance(account.ID)
	return newAccount, error
}

func recuperarDadosAccountEvent(c *fiber.Ctx) (AccountEvent, error) {
	accountEvent := new(AccountEvent)
	erro := c.BodyParser(accountEvent)
	return *accountEvent, erro
}

func jsonConvert(result interface{}) string {
	resultJSON, _ := json.Marshal(result)
	return string(resultJSON)

}
