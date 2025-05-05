package domain

import (
	"context"
	"errors"
)

var ErrPetNotFound = errors.New("repository: pet not found")

type IPetRepository interface {
	FindById(ctx context.Context, id PetId) (*Pet, error)
	FindByUserId(ctx context.Context, userId UserId) ([]*Pet, error)
	Create(ctx context.Context, pet *Pet) error
}
