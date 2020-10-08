package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Exemplo represents a model
type Exemplo struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Nome string             `json:"nome" bson:"nome"`
}

// ErrorResponse represents an error response to the user
type ErrorResponse struct {
	Msg  string `json:"msg,omitempty"`
	Erro string `json:"erro,omitempty"`
}
