package main

import (
	gmc "github.com/eduardohitek/golang-mongo-common"
	"go.mongodb.org/mongo-driver/mongo"
)

// DB represents a DB Object
type DB struct {
	Client  *mongo.Client
	DBUrl   string
	DBName  string
	DBUser  string
	DBPass  string
	DBLocal string
	Log     *LogHeimdall
}

// NewDB Return a new DB
func NewDB(DBUrl string, DBName string, DBUser string, DBPass string, DBLocal string, Log *LogHeimdall) *DB {
	return &DB{DBUrl: DBUrl, DBName: DBName, DBUser: DBUser, DBPass: DBPass, DBLocal: DBLocal, Log: Log}
}

func (db *DB) connectarDB(serviceName string) {
	erro := db.retornarClient(serviceName)
	if erro != nil {
		db.Log.Logger.Println("Erro ao se connectar com o DB", erro.Error())
		return
	}
	db.Log.Logger.Println("Connected to DB!")
}

func (db *DB) retornarClient(serviceName string) error {
	var erro error
	if db.DBLocal == "Y" {
		db.Client, erro = gmc.RetornarCliente(db.DBUrl, serviceName)
	} else {
		db.Client, erro = gmc.RetornarClienteSeguro(db.DBUrl, "admin", db.DBUser, db.DBPass, serviceName)
	}
	if erro != nil {
		db.Log.Logger.Fatalln(erro.Error())
	}
	return erro
}
