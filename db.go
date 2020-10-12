package main

import (
	"context"

	gmc "github.com/eduardohitek/golang-mongo-common"
	"go.mongodb.org/mongo-driver/bson"
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

func (db *DB) connectDB(serviceName string) {
	erro := db.returnClient(serviceName)
	if erro != nil {
		db.Log.Logger.Println("Erro ao se connectar com o DB", erro.Error())
		return
	}
	db.Log.Logger.Println("Connected to DB!")
}

func (db *DB) returnClient(serviceName string) error {
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

func (db *DB) resetDatabase() (*mongo.DeleteResult, error) {
	return gmc.Deletar(db.DBName, "accounts", db.Client, bson.M{})
}

func (db *DB) getBalance(ID string) (*Account, error) {
	account := &Account{}
	filter := bson.M{"id": ID}
	collection := db.Client.Database(db.DBName).Collection("accounts")
	result := collection.FindOne(context.TODO(), filter)
	erro := result.Decode(&account)
	if erro != nil {
		db.Log.Logger.Println("Error on retrieving the account", erro.Error())
		return nil, erro
	}
	return account, nil
}

func (db *DB) createAccount(account Account) (*mongo.InsertOneResult, error) {
	collection := db.Client.Database(db.DBName).Collection("accounts")
	result, error := collection.InsertOne(context.TODO(), account)
	return result, error
}

func (db *DB) depositAmount(account Account) error {
	filter := bson.M{"id": account.ID}
	atualizacao := bson.D{{Key: "$inc", Value: bson.M{"balance": account.Balance}}}
	collection := db.Client.Database(db.DBName).Collection("accounts")
	_, error := collection.UpdateOne(context.TODO(), filter, atualizacao)
	return error
}
