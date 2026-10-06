package usecase

import (
	"context"
	"time"

	"dashboard-fin/internal/domain"
)

type CreateFamilyInput struct {
	Name string `json:"name"`
}

type AddMemberInput struct {
	FamilyID   string  `json:"family_id"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	BaseIncome float64 `json:"base_income"`
}

type UpdateMemberInput struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	BaseIncome float64 `json:"base_income"`
}

type FamilyUseCase interface {
	CreateFamily(ctx context.Context, input CreateFamilyInput) (*domain.Family, error)
	GetFamily(ctx context.Context, id string) (*domain.Family, error)
	AddMember(ctx context.Context, input AddMemberInput) (*domain.Member, error)
	ListMembers(ctx context.Context, familyID string) ([]domain.Member, error)
	UpdateMember(ctx context.Context, input UpdateMemberInput) (*domain.Member, error)
}

type familyUseCase struct {
	repo domain.FamilyRepository
}

func NewFamilyUseCase(repo domain.FamilyRepository) FamilyUseCase {
	return &familyUseCase{repo: repo}
}

func (uc *familyUseCase) CreateFamily(ctx context.Context, input CreateFamilyInput) (*domain.Family, error) {
	if input.Name == "" {
		return nil, domain.ErrInvalidName
	}

	family := &domain.Family{
		Name:      input.Name,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.CreateFamily(ctx, family); err != nil {
		return nil, err
	}

	return family, nil
}

func (uc *familyUseCase) GetFamily(ctx context.Context, id string) (*domain.Family, error) {
	return uc.repo.FindFamilyByID(ctx, id)
}

func (uc *familyUseCase) AddMember(ctx context.Context, input AddMemberInput) (*domain.Member, error) {
	if input.Name == "" {
		return nil, domain.ErrInvalidName
	}
	if input.Role == "" {
		input.Role = "member"
	}

	member := &domain.Member{
		FamilyID:   input.FamilyID,
		Name:       input.Name,
		Role:       input.Role,
		BaseIncome: input.BaseIncome,
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (uc *familyUseCase) ListMembers(ctx context.Context, familyID string) ([]domain.Member, error) {
	return uc.repo.ListMembersByFamilyID(ctx, familyID)
}

func (uc *familyUseCase) UpdateMember(ctx context.Context, input UpdateMemberInput) (*domain.Member, error) {
	if input.ID == "" {
		return nil, domain.ErrMemberNotFound
	}

	member := &domain.Member{
		ID:         input.ID,
		Name:       input.Name,
		Role:       input.Role,
		BaseIncome: input.BaseIncome,
	}

	if err := uc.repo.UpdateMember(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}
