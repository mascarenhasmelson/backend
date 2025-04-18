package models 

import(
	"time"	
)
type Account struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Name       string    `json:"name"`
	Balance    float64   `json:"balance"`
	CreatedAt  time.Time `json:"created_at"`
}

type Account_ledger struct {
    ID         string    `json:"id" bson:"id"`
	Balance    float64   `json:"balance" bson:"balance"`
	Type       string    `json:"type" bson:"type"`
	Amount     float64 `bson:"amount" json:"amount"`
    CustomerID string    `json:"customer_id" bson:"customer_id"`
    Name       string    `json:"name" bson:"name"`
	Withdraw_balance float64 `json:"withdraw_balance" bson:"withdraw_balance"`
	Deposit_balance float64   `json:"deposit_balance" bson:"deposit_balance"`
    CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}
type PubAccount struct {
	CustomerID string  `json:"customerID"` 
	Name       string  `json:"name"`       
	Amount     float64 `json:"amount"`    
}

type Transaction struct {
	CustomerID string  `json:"customerID"` 
	Amount     float64 `json:"amount"`    
	IsDeposit  bool    `json:"isDeposit"`  
}
