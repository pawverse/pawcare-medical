package endpoint

import (
	"context"
	"errors"
	"time"

	petdomain "github.com/pawverse/pawcare-medical/internal/pet/domain"
	"github.com/pawverse/pawcare-medical/internal/record/domain"
	"github.com/pawverse/pawcare-medical/internal/record/service"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/utils"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
)

type GetByIdRequest struct {
	Id string `json:"id" validate:"required,mongodb"`
}

type GetResponse struct {
	common.EmbedError

	Id          string    `json:"id"`
	PetId       string    `json:"pet_id"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func (r GetResponse) StatusCode() int {
	var validationError validator.ValidationErrors
	if errors.As(r.Err, &validationError) {
		return 400
	}

	switch r.Err {
	case domain.ErrRecordNotFound, petdomain.ErrPetNotFound:
		return 404
	case nil:
		return 200
	default:
		return 500
	}
}

var _ http.StatusCoder = (*GetResponse)(nil)

type GetManyResponse struct {
	common.EmbedError

	Records []GetResponse `json:"records"`
}

func (r GetManyResponse) StatusCode() int {
	var validationError validator.ValidationErrors
	if errors.As(r.Err, &validationError) {
		return 400
	}

	switch r.Err {
	case nil:
		return 200
	default:
		return 500
	}
}

var _ http.StatusCoder = (*GetManyResponse)(nil)

func makeGetByPetIdEndpoint(recordService service.IRecordService) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(GetByIdRequest)
		if !ok {
			return GetManyResponse{EmbedError: common.NewEmbededError(common.ErrCastRequest)}, nil
		}

		if err := validator.New().Struct(req); err != nil {
			return GetManyResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		petId := domain.PetId(req.Id)
		records, err := recordService.GetByPetId(ctx, petId)
		if err != nil {
			return GetManyResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return GetManyResponse{
			Records: utils.Map(records, func(r *domain.Record) GetResponse {
				return GetResponse{
					Id:          string(r.Id),
					PetId:       string(r.PetId),
					Type:        string(r.RecordInfo.Type),
					Date:        r.RecordInfo.Date,
					Description: r.RecordInfo.Description,
				}
			}),
		}, nil
	}
}

func makeGetByIdEndpoint(recordService service.IRecordService) endpoint.Endpoint {
	return func(ctx context.Context, request any) (response any, err error) {
		req, ok := request.(GetByIdRequest)
		if !ok {
			return GetResponse{EmbedError: common.NewEmbededError(common.ErrCastRequest)}, nil
		}

		if err := validator.New().Struct(req); err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		id := domain.RecordId(req.Id)
		record, err := recordService.GetById(ctx, id)
		if err != nil {
			return GetResponse{EmbedError: common.NewEmbededError(err)}, nil
		}

		return GetResponse{
			Id:          string(record.Id),
			PetId:       string(record.PetId),
			Type:        string(record.RecordInfo.Type),
			Description: record.RecordInfo.Description,
			Date:        record.RecordInfo.Date,
		}, nil
	}
}
