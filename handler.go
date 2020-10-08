package main

import (
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

// Error Return an Error response to the user
func Error(c *fiber.Ctx, msg string, erro string) {
	c.Status(500).JSON(ErrorResponse{Msg: msg, Erro: erro})
	return
}
