package endpoint

import (
	"context"
	"time"

	"github.com/pawverse/pawcare-medical/internal/record/domain"
	"github.com/pawverse/pawcare-medical/internal/record/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	PetId       string    `json:"pet_id" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=vaccination treatment deworming surgery checkup other"`
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

func makeCreateEndpoint(recordService service.IRecordService) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(CreateRequest)
		if !ok {
			return GetResponse{EmbedError: common.NewEmbededError(common.ErrCastRequest)}, nil
		}

		if err := validator.New().Struct(req); err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		petId := domain.PetId(req.PetId)
		recordInfo := domain.NewRecordInfo(domain.RecordType(req.Type), req.Description, req.Date)
		record, err := recordService.Create(ctx, petId, recordInfo)
		if err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return GetResponse{
			Id:          string(record.Id),
			PetId:       string(record.PetId),
			Type:        string(record.RecordInfo.Type),
			Date:        record.RecordInfo.Date,
			Description: recordInfo.Description,
		}, nil
	}
}
