package main

//Account represents a model
type Account struct {
	ID      int32   `json:"id" bson:"id"`
	Balance float32 `json:"balance" bson:"balance"`
}

//AccountEvent represents a model
type AccountEvent struct {
	Type        string  `json:"type,omitempty"`
	Origin      string  `json:"origin,omitempty"`
	Destination string  `json:"destination,omitempty"`
	Amount      float32 `json:"amount,omitempty"`
}

type AccountEventResponse struct {
	Origin      Account `json:"origin,omitempty"`
	Destination Account `json:"destination,omitempty"`
}

// ErrorResponse represents an error response to the user
type ErrorResponse struct {
	Msg  string `json:"msg,omitempty"`
	Erro string `json:"erro,omitempty"`
}
