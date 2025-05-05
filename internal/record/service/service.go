package service

import (
	"context"

	petdomain "github.com/pawverse/pawcare-medical/internal/pet/domain"
	petservice "github.com/pawverse/pawcare-medical/internal/pet/service"
	"github.com/pawverse/pawcare-medical/internal/record/domain"
	"go.uber.org/zap"
)

type IRecordService interface {
	GetById(ctx context.Context, id domain.RecordId) (*domain.Record, error)
	GetByPetId(ctx context.Context, petId domain.PetId) ([]*domain.Record, error)
	Create(ctx context.Context, petId domain.PetId, record domain.RecordInfo) (*domain.Record, error)
}

type recordService struct {
	recordRepository domain.IRecordRepository
	petService       petservice.IPetService
}

func NewRecordService(recordRepository domain.IRecordRepository, petService petservice.IPetService, logger *zap.Logger) IRecordService {
	var svc IRecordService
	svc = &recordService{recordRepository: recordRepository, petService: petService}
	svc = newLoggingMiddleware(logger)(svc)

	return svc
}

func (svc *recordService) GetById(ctx context.Context, id domain.RecordId) (*domain.Record, error) {
	record, err := svc.recordRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := svc.petService.GetById(ctx, petdomain.PetId(record.PetId)); err != nil {
		return nil, err
	}

	return record, nil
}

func (svc *recordService) GetByPetId(ctx context.Context, petId domain.PetId) ([]*domain.Record, error) {
	if _, err := svc.petService.GetById(ctx, petdomain.PetId(petId)); err != nil {
		return nil, err
	}

	return svc.recordRepository.FindByPetId(ctx, petId)
}

func (svc *recordService) Create(ctx context.Context, petId domain.PetId, recordInfo domain.RecordInfo) (*domain.Record, error) {
	if _, err := svc.petService.GetById(ctx, petdomain.PetId(petId)); err != nil {
		return nil, err
	}

	recordAggregate := domain.NewRecord(petId, recordInfo)
	if err := svc.recordRepository.Create(ctx, recordAggregate); err != nil {
		return nil, err
	}

	return recordAggregate, nil
}
