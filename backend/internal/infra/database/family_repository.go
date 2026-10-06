package database

import (
	"context"

	"gorm.io/gorm"

	"dashboard-fin/internal/domain"
)

type familyRepository struct {
	db *gorm.DB
}

func NewFamilyRepository(db *gorm.DB) domain.FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) CreateFamily(ctx context.Context, family *domain.Family) error {
	return r.db.WithContext(ctx).Create(family).Error
}

func (r *familyRepository) FindFamilyByID(ctx context.Context, id string) (*domain.Family, error) {
	var family domain.Family
	err := r.db.WithContext(ctx).Preload("Members").First(&family, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrFamilyNotFound
		}
		return nil, err
	}
	return &family, nil
}

func (r *familyRepository) AddMember(ctx context.Context, member *domain.Member) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *familyRepository) ListMembersByFamilyID(ctx context.Context, familyID string) ([]domain.Member, error) {
	var members []domain.Member
	err := r.db.WithContext(ctx).Where("family_id = ?", familyID).Find(&members).Error
	return members, err
}
