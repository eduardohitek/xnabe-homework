package main

import (
	"strconv"

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
	id := c.Params("account_id", "0")
	ID, _ := strconv.ParseInt(id, 10, 32)
	account, erro := h.DB.getBalance(ID)
	if erro != nil {
		c.Status(404).SendString("0")
		return nil
	}
	c.JSON(account)
	return nil
}

// Error Return an Error response to the user
func Error(c *fiber.Ctx, msg string, erro string) {
	c.Status(500).JSON(ErrorResponse{Msg: msg, Erro: erro})
	return
}
