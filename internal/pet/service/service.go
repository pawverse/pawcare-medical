package service

import (
	"context"
	"fmt"

	"github.com/pawverse/pawcare-medical/internal/pet/domain"
	"go.uber.org/zap"
)

type IPetService interface {
	Create(ctx context.Context, id domain.PetId, userId domain.UserId) (*domain.Pet, error)
	GetById(ctx context.Context, id domain.PetId) (*domain.Pet, error)
}

func NewPetService(repository domain.IPetRepository, logger *zap.Logger) IPetService {
	var svc IPetService
	svc = &petService{repository: repository}
	svc = newLoggingMiddleware(logger)(svc)

	return svc
}

type petService struct {
	repository domain.IPetRepository
}

func (s *petService) Create(ctx context.Context, id domain.PetId, userId domain.UserId) (*domain.Pet, error) {
	pet := &domain.Pet{
		Id:     id,
		UserId: userId,
	}

	if err := s.repository.Create(ctx, pet); err != nil {
		return nil, fmt.Errorf("pet service: %w", err)
	}

	return pet, nil
}

func (s *petService) GetById(ctx context.Context, id domain.PetId) (*domain.Pet, error) {
	return s.repository.FindById(ctx, id)
}
