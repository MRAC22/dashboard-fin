package domain

import (
	"context"
	"errors"
	"time"
)

type TransactionType string

const (
	Income  TransactionType = "INCOME"
	Expense TransactionType = "EXPENSE"
)

var (
	ErrInvalidAmount = errors.New("o valor da transação deve ser maior que zero")
)

type Transaction struct {
	ID          string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FamilyID    string          `json:"family_id" gorm:"type:uuid;not null"`
	MemberID    *string         `json:"member_id,omitempty" gorm:"type:uuid"`
	Description string          `json:"description" gorm:"not null"`
	Amount      float64         `json:"amount" gorm:"type:numeric(10,2);not null"`
	Type        TransactionType `json:"type" gorm:"type:varchar(10);not null"`
	Category    string          `json:"category" gorm:"not null"`
	Date        time.Time       `json:"date" gorm:"type:date;not null"`
	CreatedAt   time.Time       `json:"created_at"`
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	FindByFamilyID(ctx context.Context, familyID string, memberID string) ([]Transaction, error)
	Delete(ctx context.Context, id string, familyID string) error
}