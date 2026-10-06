package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrFamilyNotFound = errors.New("família não encontrada")
	ErrMemberNotFound = errors.New("membro não encontrado")
	ErrInvalidName     = errors.New("o nome não pode ser vazio")
)

type Family struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Members   []Member  `json:"members,omitempty" gorm:"foreignKey:FamilyID"`
	CreatedAt time.Time `json:"created_at"`
}

type Member struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FamilyID  string    `json:"family_id" gorm:"index;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:'member'"` // e.g., admin, member, child
	CreatedAt time.Time `json:"created_at"`
}

type FamilyRepository interface {
	CreateFamily(ctx context.Context, family *Family) error
	FindFamilyByID(ctx context.Context, id string) (*Family, error)
	AddMember(ctx context.Context, member *Member) error
	ListMembersByFamilyID(ctx context.Context, familyID string) ([]Member, error)
}
