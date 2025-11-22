package models

import "time"

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID            string          `bson:"_id,omitempty" json:"id"`
	UserID        string          `bson:"userId" json:"userId"`
	Type          TransactionType `bson:"type" json:"type"`
	Amount        float64         `bson:"amount" json:"amount"`
	Currency      string          `bson:"currency" json:"currency"`
	CategoryID    *string         `bson:"categoryId,omitempty" json:"categoryId,omitempty"`
	SubcategoryID *string         `bson:"subcategoryId,omitempty" json:"subcategoryId,omitempty"`
	Description   *string         `bson:"description,omitempty" json:"description,omitempty"`
	Date          time.Time       `bson:"date" json:"date"`
	CreatedAt     time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time       `bson:"updatedAt" json:"updatedAt"`
}
