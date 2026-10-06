package database

import (
	"context"

	"gorm.io/gorm"

	"dashboard-fin/internal/domain"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) FindByFamilyID(ctx context.Context, familyID string, memberID string) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	query := r.db.WithContext(ctx).Where("family_id = ?", familyID)

	if memberID != "" {
		query = query.Where("member_id = ?", memberID)
	}

	err := query.Order("date desc").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) Delete(ctx context.Context, id string, familyID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND family_id = ?", id, familyID).Delete(&domain.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrTransactionNotFound
	}
	return nil
}
