package usecase

import (
	"context"
	"time"

	"dashboard-fin/internal/domain"
)

type CreateTransactionInput struct {
	FamilyID    string                 `json:"family_id"`
	MemberID    *string                `json:"member_id"`
	Description string                 `json:"description"`
	Amount      float64                `json:"amount"`
	Type        domain.TransactionType `json:"type"`
	Category    string                 `json:"category"`
	Date        time.Time              `json:"date"`
}

type TransactionUseCase interface {
	Create(ctx context.Context, input CreateTransactionInput) (*domain.Transaction, error)
	ListByFamily(ctx context.Context, familyID string, memberID string) ([]domain.Transaction, error)
}

type transactionUseCase struct {
	repo domain.TransactionRepository
}

func NewTransactionUseCase(repo domain.TransactionRepository) TransactionUseCase {
	return &transactionUseCase{repo: repo}
}

func (uc *transactionUseCase) Create(ctx context.Context, input CreateTransactionInput) (*domain.Transaction, error) {
	if input.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	tx := &domain.Transaction{
		FamilyID:    input.FamilyID,
		MemberID:    input.MemberID,
		Description: input.Description,
		Amount:      input.Amount,
		Type:        input.Type,
		Category:    input.Category,
		Date:        input.Date,
		CreatedAt:   time.Now(),
	}

	if err := uc.repo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (uc *transactionUseCase) ListByFamily(ctx context.Context, familyID string, memberID string) ([]domain.Transaction, error) {
	return uc.repo.FindByFamilyID(ctx, familyID, memberID)
}
