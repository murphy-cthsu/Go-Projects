package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID      int `json:"id"`
	FirstName string `json:"FirstName"`
	LastName string `json: "LastName"`
	Number int64 `json:"number"`
	Balance int64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName: lastName,
		Number: int64(rand.Intn(1000000)),
		Balance:0,
		CreatedAt: time.Now().UTC(),
	}
}

type CreateAccountRequest struct {	
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
}

type TransferRequest struct {
	// From int `json:"from"`
	To int `json:"to"`
	Amount int `json:"amount"`
}

