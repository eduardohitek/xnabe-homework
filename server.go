package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server represents a Server object
type Server struct {
	Fiber       *fiber.App
	H           *Handler
	DB          *DB
	Log         *LogHeimdall
	ServiceName string
}

// NewServer returns a new Server
func NewServer(serviceName string) *Server {

	app := fiber.New(fiber.Config{
		ETag:                  true,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     serviceName + " - ${time} ${header:email} ${host} ${status} ${method} ${url} - ${latency}\n",
		TimeFormat: "2006/01/02 15:04:05",
	}))
	log := NewLoggerHeimdall(serviceName)
	db := createDB(log)
	handler := NewHandler(db, log)
	return &Server{Fiber: app, H: handler, DB: db, Log: log, ServiceName: serviceName}
}

func (s *Server) initRoutes() {
	s.Fiber.Get("/", s.H.handlerRetornarTodas)
}

func (s *Server) init() {
	s.initRoutes()
	s.DB.connectarDB(s.ServiceName)
}

func (s *Server) run() {
	s.init()
	s.Log.Logger.Panic(s.Fiber.Listen(":8086"))
}

func recuperarEnvVar(nome string, valorPadrao string) string {
	ret := os.Getenv(nome)
	if ret == "" {
		ret = valorPadrao
	}
	return ret
}

func createDB(log *LogHeimdall) *DB {
	DBUrl := recuperarEnvVar("DB_URL_", "localhost")
	DBName := recuperarEnvVar("DB_NAME", "tu-dev")
	DBUser := recuperarEnvVar("DB_USER", "tactus")
	DBPass := recuperarEnvVar("DB_PASS", "pass")
	DBLocal := recuperarEnvVar("DB_LOCAL", "Y")
	return NewDB(DBUrl, DBName, DBUser, DBPass, DBLocal, log)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	NewLoggerHeimdall("Cardea Servico MS").Logger.Println("Ocorreu um erro ao processar sua requisição", err.Error())
	ctx.Status(code).JSON(fiber.Map{"msg": "Ocorreu um erro ao processar sua requisição", "erro": err.Error()})
	return nil
}
